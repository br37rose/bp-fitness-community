import React from "react";

function FitnessPlanDisplay({ weeklyFitnessPlans, label, helpText }) {
  return (
    <div className="">
      {label && <label>{label}</label>}
      {weeklyFitnessPlans.map((weeklyPlan, index) => (
        <div key={index} className="mb-5 is-family-secondary">
          <h2 className="is-size-4 mb-3 has-text-weight-bold is-underlined grey-1">
            {weeklyPlan.title}
          </h2>
          {weeklyPlan.dailyPlans.map((dailyPlan, dailyIndex) => (
            <div key={dailyIndex} className="mb-3">
              <h3 className="is-size-6 mb-2 has-text-weight-bold grey-2">
                {dailyPlan.title}
              </h3>
              <p className="is-size-6 is-family-sans-serif grey-3 has-text-weight-normal">
                {dailyPlan.planDetails}
              </p>
            </div>
          ))}
        </div>
      ))}
      {helpText !== undefined && helpText !== null && helpText !== "" && (
        <p className="help">{helpText}</p>
      )}
    </div>
  );
}

export default FitnessPlanDisplay;
