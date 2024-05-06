import React from 'react';
import Breadcrumb from '../Reusable/Breadcrumb';
import Footer from './Footer';

const Layout = ({ children, breadcrumbItems }) => {
    return (
        <div className="container is-fluid">
            <section className="section is-hidden-touch">
                {breadcrumbItems &&
                    <Breadcrumb breadcrumbItems={breadcrumbItems} />}
                {children}
            </section>
            <div className="is-hidden-desktop">
                {breadcrumbItems &&
                    <Breadcrumb breadcrumbItems={breadcrumbItems} />}
                {children}
            </div>
            <Footer />
        </div>
    );
}

export default Layout;
