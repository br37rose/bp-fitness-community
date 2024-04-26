import { Link } from "react-router-dom";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faCalendarMinus,
  faCalendarPlus,
  faEye,
} from "@fortawesome/free-solid-svg-icons";
import {
  PAGE_SIZE_OPTIONS,
  WORKOUT_STATUS_MAP,
} from "../../../Constants/FieldOptions";
import { currentUserState } from "../../../AppState";
import { useRecoilState } from "recoil";
import { patchFitnessChaleengeParticipation } from "../../../API/FitnessChallenge";

function MemberListDesktop(props) {
  const {
    listData,
    setPageSize,
    pageSize,
    previousCursors,
    onPreviousClicked,
    onNextClicked,
    onChangeStatus,
  } = props;
  const [currentUser] = useRecoilState(currentUserState);

  return (
    <div className="b-table">
      <div className="table-wrapper has-mobile-cards">
        <table className="table is-fullwidth is-striped is-hoverable is-fullwidth">
          <thead>
            <tr>
              <th>Name</th>
              <th>Users count</th>
              <th>Status</th>
              <th>CreatedAt</th>
              <th></th>
            </tr>
          </thead>
          <tbody>
            {listData &&
              listData.results &&
              listData.results.map(function (datum, i) {
                return (
                  <tr key={`desktop_${datum.id}`}>
                    <td data-label="Name">{datum.name}</td>
                    <td data-label="users-count">{datum.userIds?.length}</td>

                    <td data-label="Status">
                      {WORKOUT_STATUS_MAP[datum.status]}
                    </td>
                    <td data-label="User">{datum.createdAt || " - "}</td>

                    <td className="is-actions-cell">
                      <div className="buttons is-right">
                        <Link
                          to={`/fitness-challenge/${datum.id}`}
                          className="button is-small is-dark"
                          type="button"
                        >
                          <FontAwesomeIcon className="mdi" icon={faEye} />
                          &nbsp;View
                        </Link>
                        {datum.userIds &&
                        datum.userIds
                          .filter((u) => u === currentUser.id)
                          .map((f) => f).length > 0 ? (
                          <button
                            className="button is-small is-danger"
                            onClick={() => onChangeStatus(datum.id)}
                          >
                            <FontAwesomeIcon
                              className="mdi"
                              icon={faCalendarMinus}
                            />
                            &nbsp;Leave
                          </button>
                        ) : (
                          <button
                            className="button is-small is-success"
                            onClick={() => onChangeStatus(datum.id)}
                          >
                            <FontAwesomeIcon
                              className="mdi"
                              icon={faCalendarPlus}
                            />
                            &nbsp;Join
                          </button>
                        )}
                      </div>
                    </td>
                  </tr>
                );
              })}
          </tbody>
        </table>

        <div class="columns">
          <div class="column is-half">
            <span class="select">
              <select
                class={`input has-text-grey-light`}
                name="pageSize"
                onChange={(e) => setPageSize(parseInt(e.target.value))}
              >
                {PAGE_SIZE_OPTIONS.map(function (option, i) {
                  return (
                    <option
                      selected={pageSize === option.value}
                      value={option.value}
                    >
                      {option.label}
                    </option>
                  );
                })}
              </select>
            </span>
          </div>
          <div class="column is-half has-text-right">
            {previousCursors.length > 0 && (
              <button class="button" onClick={onPreviousClicked}>
                Previous
              </button>
            )}
            {listData.hasNextPage && (
              <>
                <button class="button" onClick={onNextClicked}>
                  Next
                </button>
              </>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}

export default MemberListDesktop;
