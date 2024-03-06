import React, { useState, useEffect } from "react";

function DataDisplayRowTags(props) {

    ////
    //// Props.
    ////

    const {
        label="Tags (Optional)",
        tags=[],
        helpText=""
    } = props;


    useEffect(() => {
        let mounted = true;

        if (mounted) {
        }

        return () => { mounted = false; }
    }, []);

    ////
    //// Component rendering.
    ////

    return (
        <div class="field pb-4">
            <label class="label">{label}</label>
            <div class="control">
                <p>
                    {tags && tags.map(function(datum, i){
                        return <span class="tag is-success mr-2 mb-2">{datum.text}</span>;
                    })}
                </p>
                {helpText !== undefined && helpText !== null && helpText !== "" && <p class="help">{helpText}</p>}
            </div>
        </div>
    );
}

export default DataDisplayRowTags;
