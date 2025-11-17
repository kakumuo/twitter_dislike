package dislikes

type DislikeRecord struct {
	OwnerId   string `bson:"ownerId" datastore:"ownerId"`
	UpdatedAt int    `bson:"updatedAt,$date" datastore:"updatedAt"`
}

type TweetRecord struct {
	OwnerId   string                   `bson:"ownerId" datastore:"ownerId"`
	TweetId   string                   `bson:"tweetId" datastore:"tweetId"`
	DislikeBy map[string]DislikeRecord `bson:"dislikedBy" datastore:"dislikedBy"`
	CreatedAt int                      `bson:"createdAt,$date" datastore:"createdAt"`
	UpdatedAt int                      `bson:"updatedAt,$date" datastore:"updatedAt"`
}

type Config struct {
	ConnectionString string `bson:"connectionString" datastore:"connectionString"`
	UserName         string `bson:"username" datastore:"username"`
	Password         string `bson:"password" datastore:"password"`
	HostName         string `bson:"hostname" datastore:"hostname"`
	Port             int    `bson:"port" datastore:"port"`
}
