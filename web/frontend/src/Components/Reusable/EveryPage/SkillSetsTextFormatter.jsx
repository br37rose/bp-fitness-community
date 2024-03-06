import React from "react";


function SkillSetsTextFormatter(props) {

    ////
    //// Props.
    ////

    const {
        skillSets=[],
    } = props;


    ////
    //// Component rendering.
    ////

    return (
        <>
            {skillSets && skillSets.map(function(datum, i){
                return <span class="tag is-success mr-2 mb-2">{datum.subCategory}</span>;
            })}
        </>
    );
}

export default SkillSetsTextFormatter;
