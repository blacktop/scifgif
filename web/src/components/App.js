import React from 'react';
import axios from 'axios';
import Header from './Header';
import SearchBar from './SearchBar';
import ImageList from './ImageList';

class App extends React.Component {
  state = { images: [] };

  onSearchSubmit = async term => {
    const response = await axios.get('/web/search', {
      params: { query: term , type: "giphy" }
    })
    // console.log('response.data', response.data)
    if (response.data) {
      this.setState({ images: response.data });
    }
  };

  render() {
    return (
      <div className="ui container" style={{ marginTop: '10px' }}>
        <Header />
        <SearchBar onSubmit={this.onSearchSubmit} />
        <ImageList images={this.state.images} />
      </div>
    );
  }
}

export default App;