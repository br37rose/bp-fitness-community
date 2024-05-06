import { useEffect, useState } from "react";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faExternalLink } from "@fortawesome/free-solid-svg-icons";
import { useLocation } from "react-router-dom";

const Footer = () => {
  const [currentDateTime, setCurrentDateTime] = useState("");
  
  const location = useLocation();
  const isindexPage = location.pathname === '/'

  useEffect(() => {
    // Function to update current date and time
    const updateDateTime = () => {
      const now = new Date();
      const options = {
        year: "numeric",
        month: "long",
        day: "numeric",
        hour: "numeric",
        minute: "numeric",
        second: "numeric",
        hour12: true,
      };
      const formattedDateTime = now.toLocaleString("en-US", options);
      setCurrentDateTime(formattedDateTime);
    };

    // Update the current date and time initially and then every minute
    updateDateTime();
    const interval = setInterval(updateDateTime, 60000); // 60000ms = 1 minute

    // Clean up the interval on component unmount
    return () => clearInterval(interval);
  }, []);

  return (
    <footer className="footer-main px-5 py-4" role="contentinfo">
      <div className="container is-fluid">
        <div className="columns is-flex-wrap-wrap">
          <div className="column is-half-desktop">
            <div className="copy-right">
              <p className="is-size-7 lightGrey" id="currentDateTime">
                {currentDateTime} in Ontario <br />
                Developed by
                <a
                  href="https://bcinnovationlabs.com"
                  rel="noopener noreferrer"
                  target="_blank"
                >
                  &nbsp;BCI Innovation Labs&nbsp;
                  <FontAwesomeIcon className="fas" icon={faExternalLink} />
                </a>
              </p>
            </div>
          </div>
          <div className="column is-half-desktop">
            <div className="is-flex-desktop is-justify-content-end">
              <div className="">
                <p className="is-size-7 mr-5 lightGrey" id="siteInfo">
                  <a
                    href="https://bp8.ca/privacy-policy/"
                    rel="noopener noreferrer"
                    target="_blank"
                  >
                    Privacy Policy&nbsp;
                    <FontAwesomeIcon className="fas" icon={faExternalLink} />
                  </a>
                  &nbsp;©2023&nbsp; BP8 Fitness Training <br />
                  1828 Blue Heron Dr Unit 26, London ON N6H 0B7
                </p>
              </div>
              <div style={{ maxWidth: "210px" }}>
                <a href="#" aria-label="F45 Training Logo">
                  <img src={isindexPage ? "/static/logo.png":"/static/footer.png"} alt="F45 Training Logo" />
                </a>
              </div>
            </div>
          </div>
        </div>
      </div>
    </footer>
  );
};

export default Footer;