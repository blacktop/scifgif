import React, { Component } from "react";

class Controls extends Component {
  constructor(props) {
    super(props);
    this.state = { selectedOption: "giphy" };
  }

  handleOptionChange(changeEvent) {
    this.setState({
      selectedOption: changeEvent.target.value
    });
  }

  handleFormSubmit(formSubmitEvent) {
    formSubmitEvent.preventDefault();

    console.log("You have selected:", this.state.selectedOption);
  }

  render() {
    return (
      <form>
        <div className="btn-group" data-toggle="buttons">
          <label className="btn btn-primary active">
            <input
              type="radio"
              value="giphy"
              checked={this.state.selectedOption === "giphy"}
              onChange={this.handleOptionChange}
            />
            giphy
          </label>
          <label className="btn btn-primary">
            <input
              type="radio"
              value="xkcd"
              checked={this.state.selectedOption === "xkcd"}
              onChange={this.handleOptionChange}
            />
            xkcd
          </label>
          <label className="btn btn-primary" disabled>
            <input
              type="radio"
              value="dilbert"
              checked={this.state.selectedOption === "dilbert"}
              onChange={this.handleOptionChange}
              disabled
            />
            dilbert
          </label>
        </div>
      </form>
    );
  }

  handleClick(size) {
    this.props.addImage(size);
  }
}

export default Controls;
