import React from "react";

function PrivacyStaticPage() {
  // For debugging purposes only.
  console.log("REACT_APP_WWW_PROTOCOL:", process.env.REACT_APP_WWW_PROTOCOL);
  console.log("REACT_APP_WWW_DOMAIN:", process.env.REACT_APP_WWW_DOMAIN);
  console.log("REACT_APP_API_PROTOCOL:", process.env.REACT_APP_API_PROTOCOL);
  console.log("REACT_APP_API_DOMAIN:", process.env.REACT_APP_API_DOMAIN);

  ////
  //// Component rendering.
  ////

  return (
    <>
      <div className="container">
        <section className="section content">
          <h1>Privacy Policy</h1>
          <p>This Privacy Policy explains how BP8 Community ("we", "our", "us") collects, uses, and discloses personal information through our fitness web application and services. By using our application, you agree to the terms of this Privacy Policy.</p>
          <h2>Information We Collect</h2>
          <h3>Biometric Data</h3>
          <p>We collect biometric data such as heart rate from wearable fitness devices that you connect to our application. This biometric data is used to provide you with personalized fitness tracking and recommendations within the BP8 Community application.</p>
          <h3>Personal Information</h3>
          <p>We may also collect personal information such as your name, email address, age, and location when you create an account or use certain features of our application.</p>
          <h3>Usage Data</h3>
          <p>We automatically collect usage data such as your IP address, browser type, device information, and interactions with our application. This helps us analyze application usage and improve our services.</p>
          <h2>How We Use Your Information</h2>
          <p>We use the personal and biometric information collected to:</p>
          <ul>
           <li>Provide and improve the BP8 Community application and services</li>
           <li>Personalize your fitness experience with tailored recommendations</li>
           <li>Communicate with you about your account or updates</li>
           <li>Analyze usage trends to enhance our offerings</li>
          </ul>
          <p>We will never sell or rent your personal or biometric data to third parties.</p>
          <h2>Data Security</h2>
          <p>We implement reasonable security measures to protect the personal and biometric information we collect from unauthorized access or disclosure. However, no method of transmission or storage is 100% secure.</p>
          <h2>Your Rights</h2>
          <p>You may access, update, or delete your account information by logging into your BP8 Community account. If you wish to have your biometric data removed from our systems, please contact us at privacy@bp8community.com.</p>
          <h2>Children's Privacy</h2>
          <p>Our application is not intended for use by children under 13. We do not knowingly collect personal information from children.</p>
          <h2>Updates</h2>
          <p>We may update this Privacy Policy from time to time to reflect changes to our practices or for other operational reasons. The updated policy will be posted on our website.</p>
          <h2>Contact Us</h2>
          <p>If you have any questions about this Privacy Policy, please contact us at privacy@bp8community.com.</p>          
        </section>
      </div>
    </>
  );
}

export default PrivacyStaticPage;
