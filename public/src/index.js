import _ from "lodash";
import React, { Component } from "react";
import ReactDOM from "react-dom";
import axios from "axios";
// import database from "database";
import SearchBar from "./components/search_bar";
import ResultDetail from "./components/result_detail";
import ResultList from "./components/result_list";
import Header from "./components/Header";
import Gallery from "./components/Gallery";

// let client = new database.Client({ host: "localhost:9200", log: "error" });
// const searchSize = 100;

class App extends Component {
  constructor(props) {
    super(props);

    this.state = {
      results: [],
      selectedResult: null
    };

    this.eSearch("happy");
  }

  eSearch(term) {
    // GET request for remote image
    axios
      .get("/web/search", {
        params: {
          query: term,
          type: "giphy"
        }
      })
      .then(response => {
        this.setState({
          results: response.data,
          selectedResult: response.data[0]
        });
      })
      .catch(error => {
        console.log(error);
      });
  }

  render() {
    const eSearch = _.debounce(term => {
      this.eSearch(term);
    }, 300);

    return (
      <div>
        <Header />
        <SearchBar onSearchTermChange={eSearch} />
        <Gallery results={this.state.results} />
      </div>
    );
  }
}

ReactDOM.render(<App />, document.querySelector(".container"));
