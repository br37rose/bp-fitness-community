import React from "react";


/**
  Example:

  // STEP 1: Example of a map.
  const PHYSICAL_ACTIVITY_SEDENTARY = 1;
  const PHYSICAL_ACTIVITY_LIGHTLY_ACTIVE = 2;
  const PHYSICAL_ACTIVITY_MODERATELY_ACTIVE = 3;
  const PHYSICAL_ACTIVITY_VERY_ACTIVE = 4;
  const PHYSICAL_ACTIVITY_MAP = {
      [PHYSICAL_ACTIVITY_SEDENTARY]: "Sedentary",
      [PHYSICAL_ACTIVITY_LIGHTLY_ACTIVE]: "Lightly Active",
      [PHYSICAL_ACTIVITY_MODERATELY_ACTIVE]: "Moderately Active",
      [PHYSICAL_ACTIVITY_VERY_ACTIVE]: "Very Active",
  };

  // ...

  // STEP 2: Create the state.
  const [datum, setDatum] = useState({physicalActivity:2});

  // ...

  // STEP 3: Use the component
  <DataDisplayRowSelectStatic
     label="My current level of physical activity is"
     selectedValue={datum.physicalActivity}
     map={PHYSICAL_ACTIVITY_MAP}
  />
 */
function DataDisplayRowSelectStatic(props) {
    const { label="", selectedValue=0, map={}, helpText="" } = props;

    console.log("selectedValue", selectedValue, "map", map) // For debugging purposes only.

    return (
        <div class="field pb-4">
            <label class="label">{label}</label>
            <div class="control">
                <p>{selectedValue && map[selectedValue]}</p>
                {helpText !== undefined && helpText !== null && helpText !== "" && <p class="help">{helpText}</p>}
            </div>
        </div>
    );
}

export default DataDisplayRowSelectStatic;
