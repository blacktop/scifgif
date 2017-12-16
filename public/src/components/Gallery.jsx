import React, { Component } from "react";
import { CopyToClipboard } from "react-copy-to-clipboard";

/* #################### */
/* ##### Gallery ###### */
/* #################### */

export default class Gallery extends Component {
  constructor(props) {
    super(props);
    this.state = {
      gallery: [],
      copied: false,
      value: ""
    };
  }
  render() {
    // const images = this.state.gallery.map((image, key) => {
    const images = this.props.results.map((image, key) => {
      return (
        <Image
          key={key}
          id={image.id}
          path={image.path}
          deleteImage={this.deleteImage.bind(this)}
        />
      );
    });
    return (
      <div>
        <ul className="grid">{images}</ul>
      </div>
    );
  }

  addImage(format) {
    const { gallery } = this.state;
    const galleryLength = gallery.length;

    const newImage = {
      id: galleryLength + 1,
      format: format,
      width: width,
      height: height
    };
    this.setState({
      gallery: gallery.concat(newImage)
    });
  }

  deleteImage(id) {
    const newState = this.state.gallery.filter(item => {
      return item.id !== id;
    });
    this.setState({
      gallery: newState
    });
  }
}

/* ####################### */
/* ##### Image Item ###### */
/* ####################### */

class Image extends Component {
  render() {
    const { id, path } = this.props;

    return (
      <li className={`grid__item grid__item--big`}>
        <CopyToClipboard
          text={`http://${window.location.hostname}:3993/${path}`}
        >
          <img
            className="grid__image"
            src={`http://${window.location.hostname}:3993/${path}`}
            alt=""
          />
        </CopyToClipboard>
        <button
          className="grid__close"
          onClick={this.handleDelete.bind(this, id)}
        >
          <span className="fas fa-trash-alt" />
        </button>
      </li>
    );
  }

  handleDelete(id) {
    this.props.deleteImage(id);
  }
}
