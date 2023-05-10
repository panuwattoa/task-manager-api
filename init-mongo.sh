mongosh -- "$MONGO_INITDB_DATABASE" <<EOF
    var rootUser = '$MONGO_INITDB_ROOT_USERNAME';
    var rootPassword = '$MONGO_INITDB_ROOT_PASSWORD';
    var admin = db.getSiblingDB('admin');
    admin.auth(rootUser, rootPassword);
    var user = '$MONGO_INITDB_USERNAME';
    var passwd = '$MONGO_INITDB_PASSWORD';

    db.createUser({
      user: user,
      pwd: passwd,
      roles: ["readWrite"]
    });

    db.createCollection("tasks");
    db.createCollection("profiles");
    db.createCollection("comments");

  db.profiles.insertMany([
    {
        "owner_id": "1234",
        "email": "panuwat_jarujareet@hotmail.com",
        "display_name": "Panuwat Jarujareet",
        "update_date": 1620000000,
        "create_date": 1620000000
    },
    {
        "owner_id": "5678",
        "email": "kondee_na@gmail.com",
        "display_name": "Kondee Na",
        "update_date": 1620000000,
        "create_date": 1620000000
    }]);

    db.profiles.createIndex({ "owner_id": 1 }, { unique: true });
    db.tasks.insertMany([
        {
            "_id": ObjectId("645b9183fcfbc11433e23ab3"),
            "topic": "The Brightness of the Sun",
            "description": "It can be said thatno one prefers a life in darkness to a life with a bright future", 
            "status": 1,
            "create_date": 1683721846,
            "owner_id": "1234",
            "archive_date": null,
            "update_date": 1683721846
        },
        {
            "_id": ObjectId("645b91a756b7bf85c7bf0fcf"),
            "topic": "Vocabulary",
            "description": "Lets start with the first tip. First practice new words by writing sentences.", 
            "status": 1,
            "create_date": 1683723368,
            "owner_id": "5678",
            "archive_date": null,
            "update_date": 1683723368
        },
        {
            "_id": ObjectId("645b93d50c49d4df72dda2dc"),
            "topic": "First, practice new words by writing sentences.",
            "description": "While reading English articles or books you can find some interesting new words there. And try to use it later when you speak or write.", 
            "status": 1,
            "create_date": 1683723394,
            "owner_id": "1234",
            "archive_date": null,
            "update_date": 1683723394
        },
        {
            "_id": ObjectId("645b93db3d5f328160153dba"),
            "topic": "synonyms",
            "description": "You can also find the synonyms to expand your vocabulary. Synonyms are the words that have similar or related meanings to other words.", 
            "status": 1,
            "create_date": 1683723413,
            "owner_id": "5678",
            "archive_date": null,
            "update_date": 1683723413
        },
        {
            "_id": ObjectId("645b93e2deb13869ae24d96d"),
            "topic": "examples",
            "description": "Here are examples of synonyms Beautiful and pretty are synonyms. Also, Happy and joyful are synonyms.", 
            "status": 1,
            "create_date": 1683723418,
            "owner_id": "1234",
            "archive_date": null,
            "update_date": 1683723418
        },
        {
            "_id": ObjectId("645b93e74592a60e3d72368f"),
            "topic": "attract",
            "description": "attract is a verb attraction is a noun and attractive is an adjective. Let's see the example of using these words", 
            "status": 1,
            "create_date": 1683725446,
            "owner_id": "5678",
            "archive_date": 1683723423,
            "update_date": 1683723423
        }
    ]);

      db.comments.insertMany([
        {
            "topic_id": "645b9183fcfbc11433e23ab3",
            "owner_id": "1234",
            "content": "I agree with you",
            "create_date": 1683722889,
            "update_date": 1683722889
        },
        {
            "topic_id": "645b9183fcfbc11433e23ab3",
            "owner_id": "5678",
            "content": "I agree with you too",
            "create_date": 1683722889,
            "update_date": 1683722889
        }
        ]);
        db.comments.createIndex({ "topic_id": 1 });

EOF
