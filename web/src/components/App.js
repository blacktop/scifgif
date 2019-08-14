import React from 'react';
import axios from 'axios';
import Header from './Header';
import SearchBar from './SearchBar';
import Filter from './Filter';
import ImageList from './ImageList';

class App extends React.Component {
  state = { images: [], type: "giphy" };

  onFilterChange = newType => {
    this.setState({ type: newType });
  }

  onSearchSubmit = async term => {
    try {
      const response = await axios.get('/web/search', {
        params: { query: term, type: this.state.type }
      })
      // console.log('response :', response);
      this.setState({ images: response.data });
    } catch (error) {
      console.log(error.response);
    }
  };

  render() {
    return (
      <div className="ui container" style={{ marginTop: '10px' }}>
        <Header />
        <Filter onChange={this.onFilterChange} type={this.state.type} />
        <SearchBar type={this.state.type} onSubmit={this.onSearchSubmit} />
        <ImageList images={this.state.images} />
      </div>
    );
  }
}

export default App;