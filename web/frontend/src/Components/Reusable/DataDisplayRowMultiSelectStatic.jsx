import React, { useState, useEffect } from "react";


/**
  Example:

  // STEP 1: Example of a map.
  const HOME_GYM_EQUIPMENT_MAP = {
      2: "Bench/Boxes/Floor Mat",
      3: "Free Weights",
      4: "Barbell",
      5: "Cables",
      6: "Rower",
      7: "Stationary Bike",
      8: "Treadmill",
      9: "Resistant Bands",
      10: "Skipping Ropex",
      11: "Pull Up Bar",
      12: "Kettle Bells",
  };

  // ...

  // STEP 2: Create the state.
  const [datum, setDatum] = useState({homeGymEquipment:[2,3]});

  // ...

  // STEP 3: Use the component
  <DataDisplayRowMultiSelectStatic
      label="Please select all the home gym equipment that you have (Optional"
      selectedValues={datum.homeGymEquipment}
      map={HOME_GYM_EQUIPMENT_MAP}
  />
 */
function DataDisplayRowMultiSelectStatic(props) {

    ////
    //// Props.
    ////

    const {
        label="",
        selectedValues=[],
        map=[],
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

    console.log("selectedValues:", selectedValues);

    //TODO: IMPLEMENT
    return (
        <div class="field pb-4">
            <label class="label">{label}</label>
            <div class="control">
                <p>
                    {selectedValues && selectedValues.map(function(datum, i){
                        return <span class="tag is-success mr-2 mb-2">{map[datum]}</span>;
                    })}
                </p>
                {helpText !== undefined && helpText !== null && helpText !== "" && <p class="help">{helpText}</p>}
            </div>
        </div>
    );
}

export default DataDisplayRowMultiSelectStatic;
