import PropTypes from "prop-types";

const YouTubeVideo = ({ videoId, parseForVideoId = true }) => {
  const getYouTubeVideoId = (url) => {
    const match =
      url.match(/[?&]v=([^&]+)/) ||
      url.match(/(?:youtu\.be\/|\/embed\/|\/\?v=|&v=)([^\/?\&]+)/);
    return match && match[1];
  };

  if (parseForVideoId) {
    videoId = getYouTubeVideoId(videoId);
  }

  return (
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
};

YouTubeVideo.propTypes = {
  embedId: PropTypes.string.isRequired,
};

export default YouTubeVideo;
