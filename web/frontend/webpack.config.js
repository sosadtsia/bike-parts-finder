const path = require('path');

module.exports = {
  devServer: {
    setupMiddlewares: (middlewares, devServer) => {
      // Your custom middleware setup here
      return middlewares;
    }
  }
};
