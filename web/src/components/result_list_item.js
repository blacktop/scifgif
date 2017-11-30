import React from "react";

const ResultListItem = ({ result, onResultSelect }) => {
  const playName = result._source.play_name;
  const textEntry = result._source.text_entry;

  return (
    <li onClick={() => onResultSelect(result)} className="list-group-item">
      <div className="result-list media">
        <div className="media-body">
          <div className="media-heading">
            <h5 className={"mt-0"}>{playName}</h5>
            {textEntry}
          </div>
        </div>
      </div>
    </li>
  );
};

export default ResultListItem;
