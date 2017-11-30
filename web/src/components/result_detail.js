"use strict";
import React from "react";

const ResultDetail = ({ result }) => {
  if (!result) {
    return <div>Loading...</div>;
  }

  const playName = result._source.play_name;
  const speechNumber = result._source.speech_number;
  const lineNumber = result._source.line_number;
  const speaker = result._source.speaker;

  return (
    <div className="result-detail col-md-4">
      <div className="details">
        <div>
          <h2>
            Play: <code>{playName}</code>
          </h2>
        </div>
        <div>
          <h4>
            Line: <code>{lineNumber}</code>
          </h4>
        </div>
        <div>
          <h4>
            Speech: <code>{speechNumber}</code>
          </h4>
        </div>
        <div>
          <h4>
            Speaker: <code>{speaker}</code>
          </h4>
        </div>
      </div>
    </div>
  );
};

export default ResultDetail;
