conn = new Mongo("localhost:27017");
db = conn.getDB("discontdealer");
db.createCollection("clothes");