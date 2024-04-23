// import React, { useState } from "react";
// import PageLoadingContent from "./PageLoadingContent"; // Assuming PageLoadingContent is a component for displaying loading message

// const YouTubeVideo = ({ videoId, width, height }) => {
//   const [loaded, setLoaded] = useState(false);

//   const handleLoad = () => {
//     setLoaded(true);
//   };

//   return (
//     <div>
//       {!loaded && <PageLoadingContent displayMessage="Loading video" />}
//       <iframe
//         src={`https://www.youtube.com/embed/${videoId}`}
//         allow="autoplay; encrypted-media"
//         allowFullScreen
//         title="video"
//         onLoad={handleLoad}
//         style={{ width: width, height: height }}
//       />
//     </div>
//   );
// };

// export default YouTubeVideo;

import React from "react";
import PropTypes from "prop-types";

const YoutubeEmbed = ({ videoId }) => (
  <div className="video-responsive">
    <iframe
      width="853"
      height="480"
      src={`https://www.youtube.com/embed/${videoId}`}
      frameBorder="0"
      allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture"
      allowFullScreen
      title="Embedded youtube"
    />
  </div>
);

YoutubeEmbed.propTypes = {
  embedId: PropTypes.string.isRequired,
};

export default YoutubeEmbed;
