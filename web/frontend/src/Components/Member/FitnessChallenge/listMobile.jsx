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

/*
Display for both tablet and mobile.
*/
function MemberListMobile(props) {
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
    <>
      {listData &&
        listData.results &&
        listData.results.map(function (datum, i) {
          return (
            <div class="pb-2" key={`mobile_tablet_${datum.id}`}>
              <strong>Name:</strong>&nbsp;{datum.name}
              <br />
              <br />
              <strong>Users count:</strong>&nbsp;
              {datum.userIds?.length}
              <br />
              <br />
              <strong>Status:</strong>&nbsp;
              {WORKOUT_STATUS_MAP[datum.status]}
              <br />
              <br />
              <strong>created At:</strong>&nbsp;
              {datum.createdAt || " - "}
              <br />
              <br />
              {/* Tablet only */}
              <div class="is-hidden-mobile" key={`tablet_${datum.id}`}>
                <div className="buttons is-right">
                  <Link
                    to={`/fitness-challenge/${datum.id}`}
                    class="button is-small is-dark"
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
                      <FontAwesomeIcon className="mdi" icon={faCalendarMinus} />
                      &nbsp;Leave
                    </button>
                  ) : (
                    <button
                      className="button is-small is-success"
                      onClick={() => onChangeStatus(datum.id)}
                    >
                      <FontAwesomeIcon className="mdi" icon={faCalendarPlus} />
                      &nbsp;Join
                    </button>
                  )}
                </div>
              </div>
              {/* Mobile only */}
              <div class="is-hidden-tablet" key={`mobile_${datum.id}`}>
                <div class="columns is-mobile">
                  <div class="column">
                    <Link
                      to={`/fitness-challenge/${datum.id}`}
                      class="button is-small is-dark is-fullwidth"
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
                        className="mt-2 button is-small is-danger is-fullwidth"
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
                        className=" mt-2 button is-small is-success is-fullwidth"
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
                </div>
              </div>
            </div>
          );
        })}

      <div class="columns pt-4 is-mobile">
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
    </>
  );
}

export default MemberListMobile;
