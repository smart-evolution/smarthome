// @flow
import _ from 'lodash';
import { put, call } from 'redux-saga/effects';
import * as alertsActions from 'client/models/alerts/actions';
import * as alertsConstants from 'client/models/alerts/constants';
import * as actions from './actions';

function getData(agentID) {
  return fetch(`/api/agents/${agentID}/edit`)
    .then(response => {
      if (!response.ok) {
        throw new Error(`Fetching data error: ${response.statusText}`);
      }

      if (response.status === 204) {
        throw new Error('No data available');
      }

      return response.json();
    })
    .catch(e => e);
}

export function* fetchData({ agentID }: { agentID: string }): Iterable<any> {
  const data = yield call(getData, agentID);

  if (!_.isEmpty(data)) {
    const { temperature } = data;
    yield put(actions.updateTemperature(agentID, temperature));
  } else {
    yield put(alertsActions.addAlert('Fetching agent config failed', alertsConstants.ALERT_TYPE_ERROR));
  }
}

function callUpdateData(agentID, data) {
  return fetch(`/api/agents/${agentID}/edit`, {
    method: 'POST',
    body: JSON.stringify(data),
  })
    .then(response => response.json())
    .catch(() => 'Updating agent config failed');
}

export function* updateData({ agentID, data }: { agentID: string, data: Object }): Iterable<any> {
  const resp = yield call(callUpdateData, agentID, data);

  if (!_.isEmpty(resp)) {
    yield put(
      alertsActions.addAlert('Updated agent config successfuly', alertsConstants.ALERT_TYPE_INFO)
    );
  } else {
    yield put(alertsActions.addAlert('Updating agent config failed', alertsConstants.ALERT_TYPE_ERROR));
  }
}
