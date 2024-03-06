import React from 'react';
import PropTypes from 'prop-types';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faRefresh, faFilter, faPlus } from '@fortawesome/free-solid-svg-icons';
import { Link } from 'react-router-dom';

const TableHeader = ({
  headerTitle,
  headerIcon,
  fetchList,
  newMember,
  refresh,
  filter,
  showFilter,
  setShowFilter,
}) => {
  return (
    <div className="columns is-mobile">
      <div className="column">
        <h1 className="title is-4">
          <FontAwesomeIcon className="fas" icon={headerIcon} />
          &nbsp;{headerTitle}
        </h1>
      </div>
      <div className="column has-text-right">
        {refresh && 
          <button
            onClick={fetchList}
            className="button is-small is-info"
            type="button"
          >
            <FontAwesomeIcon className="mdi" icon={faRefresh} />
          </button>
        }
        &nbsp;
        {filter && 
          <button
            onClick={(e) => setShowFilter(!showFilter)}
            className="button is-small is-success"
            type="button"
          >
            <FontAwesomeIcon className="mdi" icon={faFilter} />
            &nbsp;Filter
          </button>
        }
        &nbsp;
        {newMember &&
          <Link
            to={`/admin/members/add`}
            className="button is-small is-primary"
            type="button"
          >
            <FontAwesomeIcon className="mdi" icon={faPlus} />
            &nbsp;New Member
          </Link>
        }
      </div>
    </div>
  );
};

TableHeader.propTypes = {
  fetchList: PropTypes.func.isRequired,
  currentCursor: PropTypes.string.isRequired,
  pageSize: PropTypes.number.isRequired,
  actualSearchText: PropTypes.string.isRequired,
  branchID: PropTypes.string.isRequired,
  showFilter: PropTypes.bool.isRequired,
  setShowFilter: PropTypes.func.isRequired,
};

export default TableHeader;
