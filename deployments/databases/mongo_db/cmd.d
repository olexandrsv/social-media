db = db.getSiblingDB('social-media');

db.adminCommand({
  setParameter: 1,
  logLevel: 0
})

db.createCollection('comments');
db.createCollection('messages');
db.createCollection('posts');
