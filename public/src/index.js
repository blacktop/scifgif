import _ from "lodash";
import React, { Component } from "react";
import ReactDOM from "react-dom";
import axios from "axios";
import elasticsearch from "elasticsearch";
import SearchBar from "./components/search_bar";
import ResultDetail from "./components/result_detail";
import ResultList from "./components/result_list";
import Gallery from "./components/Gallery";

let client = new elasticsearch.Client({ host: "localhost:9200", log: "error" });
const searchSize = 100;

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
    axios.get('/web/search', {
      params: {
        query: term,
        type: "giphy",
      }
    })
    .then((response) => {
      this.setState({
                results: response.data,
                selectedResult: response.data[0]
              });
    })
    .catch((error) => {
      console.log(error);
    });
    // // pin client
    // client.ping(
    //   {
    //     requestTimeout: 10000
    //   },
    //   function(error) {
    //     if (error) {
    //       console.error("elasticsearch cluster is down!");
    //     } else {
    //       console.log("successfully connected to elasticsearch cluster");
    //     }
    //   }
    // );
    // // search for term
    // client
    //   .search({
    //     index: "scifgif",
    //     type: "giphy",
    //     q: term.replace(/\W/g, ""),
    //     size: searchSize
    //   })
    //   .then(
    //     body => {
    //       let esResults = body.hits.hits;
    //       this.setState({
    //         results: esResults,
    //         selectedResult: esResults[0]
    //       });
    //     },
    //     error => {
    //       console.trace(error.message);
    //     }
    //   );
  }

  render() {
    const eSearch = _.debounce(term => {
      this.eSearch(term);
    }, 300);

    return (
      <div>
        <SearchBar onSearchTermChange={eSearch} />
        {/*  <ResultDetail result={this.state.selectedResult} />
        <ResultList
          onResultSelect={selectedResult => this.setState({ selectedResult })}
          results={this.state.results}
        /> */}
        <Gallery results={this.state.results} />
      </div>
    );
  }
}

ReactDOM.render(<App />, document.querySelector(".container"));
