const jwt = require('jsonwebtoken');

const authenticateSocket = (socket, next) => {
  const token = socket.handshake.auth.token;
  
  if (!token) {
    return next(new Error('Authentication error'));
  }
  
  try {
    const decoded = jwt.verify(token, process.env.JWT_SECRET);
    socket.user = {
      id: decoded.id,
      username: decoded.username
    };
    next();
  } catch (err) {
    next(new Error('Authentication error'));
  }
};

module.exports = { authenticateSocket };