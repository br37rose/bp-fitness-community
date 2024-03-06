import { useEffect } from "react";

// Special thanks to: https://stackoverflow.com/a/72495965

const RedirectURL = ({ url }) => {
  useEffect(() => {
    window.location.href = url;
  }, [url]);

  return <h5>Redirecting...</h5>;
};

export default RedirectURL;
