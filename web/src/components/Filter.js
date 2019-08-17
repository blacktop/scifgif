import './Filter.css';
import React from 'react';

class Filter extends React.Component {
    state = { type: this.props.type };

    handleChange = event => {
        this.setState({ type: event.target.value })
        this.props.onChange(event.target.value);
    };

    render() {
        const { type } = this.state
        return (
            <div className="container type-filter">
                <div className="row justify-content-center" onChange={event => this.handleChange(event)}>
                    <div class="form-check form-check-inline">
                        <input className="form-check-input" type="radio" checked={type == "giphy"} value="giphy" name="giphy" id="giphyRadio" />
                        <label className="form-check-label" for="giphyRadio">giphy</label>
                    </div>
                    <div class="form-check form-check-inline">
                        <input className="form-check-input" type="radio" checked={type == "xkcd"} value="xkcd" name="xkcd" id="xkcdRadio" />
                        <label className="form-check-label" for="xkcdRadio">xkcd</label>
                    </div>
                    <div class="form-check form-check-inline">
                        <input className="form-check-input" type="radio" checked={type == "dilbert"} value="dilbert" name="dilbert" id="dilbertRadio" />
                        <label className="form-check-label" for="dilbertRadio">dilbert</label>
                    </div>
                </div>
            </div>
        );
    }
}

export default Filter;
