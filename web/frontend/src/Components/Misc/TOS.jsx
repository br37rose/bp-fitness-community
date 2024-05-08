import React from "react";

function TermsOfServiceStaticPage() {
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
          <h1>Terms and Conditions</h1>
          <p>Welcome to BP8 Community, a fitness platform operated by BP8 Community Inc. ("Company", "we", "our", or "us"). By accessing or using our web application and services (collectively, the "Platform"), you agree to be bound by these Terms and Conditions and our Privacy Policy, which is incorporated herein by reference.</p>

          <h2>Accounts and Registration</h2>
          <p>To access certain features of the Platform, you may need to create an account. You agree to provide accurate and complete information when creating your account. You are solely responsible for any activity under your account and must keep your login credentials secure.</p>

          <h2>Biometric Data Collection</h2>
          <p>The Platform allows you to connect wearable fitness devices to track and display your biometric data, such as heart rate. By using this feature, you consent to our collection, processing, and use of your biometric data as outlined in our Privacy Policy.</p>

          <h2>User Conduct</h2>
          <p>You agree to use the Platform only for lawful purposes and in a manner that does not infringe upon the rights of others or restrict or inhibit their use of the Platform. Prohibited conduct includes:</p>
          <ul>
             <li>Uploading or transmitting any content that is unlawful, defamatory, hateful, or infringes on intellectual property rights.</li>
             <li>Attempting to gain unauthorized access to the Platform or other accounts.</li>
             <li>Interfering with the Platform's operation or security measures.</li>
          </ul>
          <p>We reserve the right to suspend or terminate your account for any violation of these Terms.</p>

          <h2>Intellectual Property</h2>
          <p>The Platform and all its content, features, and functionality are owned by the Company and protected by intellectual property laws. You may not copy, modify, distribute, or create derivative works based on the Platform without our prior written consent.</p>

          <h2>Third-Party Links and Services</h2>
          <p>The Platform may contain links to third-party websites or services. We are not responsible for the content, accuracy, or policies of any third-party sites or services linked to or integrated with the Platform.</p>

          <h2>Disclaimer of Warranties</h2>
          <p>The Platform is provided on an "as is" and "as available" basis without warranties of any kind. We do not warrant that the Platform will be uninterrupted, secure, or error-free.</p>

          <h2>Limitation of Liability</h2>
          <p>To the maximum extent permitted by law, the Company shall not be liable for any indirect, incidental, special, or consequential damages arising out of or related to your use of the Platform.</p>

          <h2>Governing Law</h2>
          <p>These Terms shall be governed by and construed in accordance with the laws of [Governing State/Country]. Any disputes shall be resolved exclusively in the courts located in [Governing City, State/Country].</p>

          <h2>Changes to Terms</h2>
          <p>We reserve the right to modify these Terms at any time. Your continued use of the Platform after any changes constitutes your acceptance of the new Terms.</p>

          <h2>Contact Us</h2>
          <p>If you have any questions about these Terms, please contact us at [contact@bp8community.com].</p>


          <p><b>Effective Date: [May 8, 2024]</b></p>
        </section>
      </div>
    </>
  );
}

export default TermsOfServiceStaticPage;
