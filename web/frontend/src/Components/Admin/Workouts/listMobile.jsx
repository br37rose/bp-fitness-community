import { Link } from "react-router-dom";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faEye } from "@fortawesome/free-solid-svg-icons";
import {
  PAGE_SIZE_OPTIONS,
  WORKOUT_STATUS_MAP,
} from "../../../Constants/FieldOptions";

/*
Display for both tablet and mobile.
*/
function AdminWorkoutListMobile(props) {
  const {
    listData,
    setPageSize,
    pageSize,
    previousCursors,
    onPreviousClicked,
    onNextClicked,
    onSelectMemberForDeletion,
  } = props;
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
              <strong>Visibility:</strong>&nbsp;
              {datum.visibility == 1 ? "Visibile to All" : "Perosonal"}
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
                    to={`/admin/workouts/${datum.id}`}
                    class="button is-small is-dark"
                    type="button"
                  >
                    <FontAwesomeIcon className="mdi" icon={faEye} />
                    &nbsp;View
                  </Link>
                  {/*
                            <Link to={`/admin/member/${datum.id}/update`} class="button is-small is-warning" type="button">
                                <FontAwesomeIcon className="mdi" icon={faPencil} />&nbsp;Edit
                            </Link>
                            <button onClick={(e, ses) => onSelectMemberForDeletion(e, datum)} class="button is-small is-danger" type="button">
                                <FontAwesomeIcon className="mdi" icon={faTrashCan} />&nbsp;Delete
                            </button>
                            */}
                </div>
              </div>
              {/* Mobile only */}
              <div class="is-hidden-tablet" key={`mobile_${datum.id}`}>
                <div class="columns is-mobile">
                  <div class="column">
                    <Link
                      to={`/admin/workouts/${datum.id}`}
                      class="button is-small is-dark is-fullwidth"
                      type="button"
                    >
                      <FontAwesomeIcon className="mdi" icon={faEye} />
                      &nbsp;View
                    </Link>
                  </div>
                  {/*
                            <div class="column">
                                <Link to={`/admin/member/${datum.id}/update`} class="button is-small is-warning is-fullwidth" type="button">
                                    <FontAwesomeIcon className="mdi" icon={faPencil} />&nbsp;Edit
                                </Link>
                            </div>
                            <div class="column">
                                <button onClick={(e, ses) => onSelectMemberForDeletion(e, datum)} class="button is-small is-danger is-fullwidth" type="button">
                                    <FontAwesomeIcon className="mdi" icon={faTrashCan} />&nbsp;Delete
                                </button>
                            </div>
                            */}
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

export default AdminWorkoutListMobile;
