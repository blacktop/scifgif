import './ImageCard.css';
import React from 'react';
import CopyToClipboard from 'react-copy-to-clipboard';

class ImageCard extends React.Component {
  constructor(props) {
    super(props);

    this.state = { spans: 0 };

    this.imageRef = React.createRef();
  }

  componentDidMount() {
    this.imageRef.current.addEventListener('load', this.setSpans);
  }

  setSpans = () => {
    const height = this.imageRef.current.clientHeight;

    const spans = Math.ceil(height / 10);

    this.setState({ spans });
  };

  render() {
    const { title, path } = this.props.image;

    return (
      <div className="image-card" style={{ gridRowEnd: `span ${this.state.spans}` }}>
        <CopyToClipboard text={`http://${window.location.hostname}:3993/${path}`} >
          <img ref={this.imageRef} alt={title} src={`http://${window.location.hostname}:3993/${path}`}/>
        </CopyToClipboard>
      </div>
    );
  }
}

export default ImageCard;
