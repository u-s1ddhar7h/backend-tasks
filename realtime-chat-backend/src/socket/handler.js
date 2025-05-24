const { authenticateSocket } = require('./auth');
const { createRoom, joinRoom, leaveRoom } = require('../utils/roomManager');

module.exports = (io) => {
  io.use(authenticateSocket);
  
  io.on('connection', (socket) => {
    console.log(`User connected: ${socket.id}`);
    
    // Room management
    socket.on('create_room', (roomName, callback) => {
      const result = createRoom(roomName, socket.id);
      callback(result);
    });
    
    socket.on('join_room', (roomId, callback) => {
      const result = joinRoom(roomId, socket.id);
      if (result.success) {
        socket.join(roomId);
      }
      callback(result);
    });
    
    socket.on('leave_room', (roomId) => {
      leaveRoom(roomId, socket.id);
      socket.leave(roomId);
    });
    
    // Messaging
    socket.on('send_message', ({ roomId, message }) => {
      socket.to(roomId).emit('receive_message', {
        sender: socket.user.username,
        message,
        timestamp: new Date()
      });
    });
    
    socket.on('disconnect', () => {
      console.log(`User disconnected: ${socket.id}`);
    });
  });
};