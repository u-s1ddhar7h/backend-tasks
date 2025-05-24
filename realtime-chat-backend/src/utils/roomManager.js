const rooms = new Map();

const createRoom = (roomName, creatorId) => {
  const roomId = generateRoomId();
  rooms.set(roomId, {
    name: roomName,
    creator: creatorId,
    participants: new Set([creatorId]),
    createdAt: new Date()
  });
  return { success: true, roomId };
};

const joinRoom = (roomId, userId) => {
  if (!rooms.has(roomId)) {
    return { success: false, error: 'Room not found' };
  }
  
  const room = rooms.get(roomId);
  room.participants.add(userId);
  
  return { success: true, roomName: room.name };
};

const leaveRoom = (roomId, userId) => {
  if (rooms.has(roomId)) {
    const room = rooms.get(roomId);
    room.participants.delete(userId);
    
    // Delete room if empty
    if (room.participants.size === 0) {
      rooms.delete(roomId);
    }
  }
};

const generateRoomId = () => {
  return Math.random().toString(36).substring(2, 10);
};

module.exports = {
  createRoom,
  joinRoom,
  leaveRoom,
  getRooms: () => Array.from(rooms.entries())
};