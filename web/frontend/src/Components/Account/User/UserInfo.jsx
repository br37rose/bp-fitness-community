import { faAddressBook, faArrowLeft, faChartPie, faContactCard, faIdCard, faPencil } from '@fortawesome/free-solid-svg-icons';
import React from 'react';
import FormTextRow from "../../Reusable/FormTextRow";
import FormTextYesNoRow from '../../Reusable/FormTextYesNoRow';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { Link } from 'react-router-dom';

const UserInfo = ({
    firstName,
    lastName,
    email,
    phone,
    country,
    region,
    city,
    addressLine1,
    addressLine2,
    postalCode,
    agreePromotionsEmail
}) => {
    return (
        <div>
            <div className='columns'>
                <div className='column is-half'>
                    <p class="title is-6"><FontAwesomeIcon className="fas" icon={faIdCard} />&nbsp;Full Name</p>
                    <FormTextRow
                        label="First Name"
                        value={firstName}
                    />
                    <FormTextRow
                        label="Last Name"
                        value={lastName}
                    />
                </div>
                <div className='column is-half'>
                    <p class="title is-6"><FontAwesomeIcon className="fas" icon={faContactCard} />&nbsp;Contact Information</p>

                    <FormTextRow
                        label="Email"
                        value={email}
                    />
                    <FormTextRow
                        label="Phone"
                        value={phone}
                    />
                </div>
            </div>

            <div className='columns'>
                <div className='column is-half'>
                    <p class="title is-6"><FontAwesomeIcon className="fas" icon={faAddressBook} />&nbsp;Address</p>

                    <FormTextRow
                        label="Country"
                        value={country}
                    />
                    <FormTextRow
                        label="Province/Territory"
                        value={region}
                    />
                    <FormTextRow
                        label="City"
                        value={city}
                    />
                    <FormTextRow
                        label="Address Line 1"
                        value={addressLine1}
                    />
                    <FormTextRow
                        label="Address Line 2"
                        value={addressLine2}
                    />
                    <FormTextRow
                        label="Postal Code"
                        value={postalCode}
                    />
                </div>
                <div className='column is-half'>
                    <p class="title is-6"><FontAwesomeIcon className="fas" icon={faChartPie} />&nbsp;Metrics</p>

                    <FormTextYesNoRow
                        label="I agree to receive electronic updates from my local gym"
                        value={agreePromotionsEmail}
                    />
                </div>
            </div>

            <div class="columns">
                <div class="column is-half">
                    <Link class="button is-medium is-fullwidth-mobile" to={"/dashboard"}><FontAwesomeIcon className="fas" icon={faArrowLeft} />&nbsp;Back to Dashboard</Link>
                </div>
                <div class="column is-half has-text-right">
                    <Link to={"/account/update"} class="button is-medium is-primary is-hidden-touch"><FontAwesomeIcon className="fas" icon={faPencil} />&nbsp;Edit</Link>
                    <Link to={"/account/update"} class="button is-medium is-primary is-fullwidth is-hidden-desktop"><FontAwesomeIcon className="fas" icon={faPencil} />&nbsp;Edit</Link>
                </div>
            </div>
        </div>
    );
};

export default UserInfo;
