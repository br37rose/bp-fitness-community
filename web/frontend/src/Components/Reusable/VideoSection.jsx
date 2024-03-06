import React from "react";
import { Link } from "react-router-dom";

const VideoSection = ({ title, items, categoryId, videoCollectionId }) => (
    <>
        <div className="columns is-multiline mt-5" role="region" aria-label="Video Section">
            <div className="column">
                <div className="is-flex is-justify-content-space-between">
                    <h3 className="is-size-4 has-text-weight-bold">{title}</h3>
                    <Link
                        to={`/video-category/${categoryId}/video-collections`}
                        className="has-text-grey-light has-text-weight-semibold"
                        aria-label="See All Videos">
                        See All
                    </Link>
                </div>
            </div>
        </div>
        <hr className="mt-0 mb-3" aria-hidden="true" />
        <div className="columns is-multiline mb-0" role="list">
            {items.map((item, idx) => (
                <div className="column mb-0 is-3" key={idx} role="listitem">
                    <Link
                        to={`/video-collection/${categoryId}/video-content/${videoCollectionId}`}
                        className="border-radius has-text-black"
                        aria-label={`Video: ${item.text}`}>
                        <img
                            className="border-radius m-w-100"
                            src={item.src}
                            alt={item.name}
                            role="presentation"
                        />
                        <h5 className="is-size-5 has-text-weight-semibold">{item.text}</h5>
                    </Link>
                </div>
            ))}
        </div>
    </>
);

export default VideoSection;
