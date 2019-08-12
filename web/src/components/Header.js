import React from 'react';

const Header = () => {
  return (
    <header className="header">
      <h1 className="header__title">
        scif[ <em className="text-primary">gif</em> ] - image search
      </h1>
      <p className="header__intro">
        Type <code>keywords</code> to filter on and then click the image to copy
        it's URL to your clipboard.
      </p>
    </header>
  );
};

export default Header;