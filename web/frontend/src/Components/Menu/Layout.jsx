import React from 'react';
import Breadcrumb from '../Reusable/Breadcrumb';

const Layout = ({ children, breadcrumbItems }) => {
    return (
        <div className="container is-fluid">
            <section className="section">
                {breadcrumbItems &&
                    <Breadcrumb breadcrumbItems={breadcrumbItems} />}
                {children}
            </section>
        </div>
    );
}

export default Layout;
