package social

import (
  "github.com/ChimeraCoder/anaconda"
  "os"
  "log"
  "strconv"
  "net/url"
)

type Twitter interface {
  GetSelfFriendIds() ([]int64, error)
  GetFriendIds(userId int64) ([]int64, error)

  GetSelfFollowerIds() ([]int64, error)
  GetFollowerIds(userId int64) ([]int64, error)

  GetUsersShowById(userId int64) (anaconda.User, error)

  Unfollow(userId int64) error
  Follow(userId int64) error
}

type TwitterConnection struct {
  api *anaconda.TwitterApi
}

func NewTwitter() TwitterConnection {

  t := TwitterConnection{}

  log.Printf("Initializing Anaconda endpoint")
	anaconda.SetConsumerKey(os.Getenv("TWITTER_CONSUMER_KEY"))
	anaconda.SetConsumerSecret(os.Getenv("TWITTER_CONSUMER_SECRET"))
	t.api = anaconda.NewTwitterApi(os.Getenv("TWITTER_TOKEN"), os.Getenv("TWITTER_TOKEN_SECRET"))
	log.Printf("Finished Initializing")

  return t
}

func (t TwitterConnection) Follow(userId int64) error {
  log.Printf("Twitter following %d", userId)
  _, err := t.api.FollowUserId(userId, nil)
  log.Printf("Finished Twitter following %d", userId)
  return err
}

func (t TwitterConnection) Unfollow(userId int64) error {
  log.Printf("Twitter unfollowing %d", userId)
  _, err := t.api.UnfollowUserId(userId)
  log.Printf("Finished Twitter unfollowing %d", userId)
  return err
}

func (t TwitterConnection) GetSelfFriendIds() ([]int64, error) {
  userId, err := getCurrentUserId()
  if err != nil {
    return nil, err
  }
  return t.GetFriendIds(userId)
}

func (t TwitterConnection) GetFriendIds(userId int64) ([]int64, error) {
  strUserId := strconv.FormatInt(userId, 10)

  v := url.Values{}
  v.Set("user_id", strUserId)
  v.Set("count", "5000")

  c, err := t.api.GetFriendsIds(v)
  if err != nil {
    return nil, err
  }

  return c.Ids, nil
}

func (t TwitterConnection) GetSelfFollowerIds() ([]int64, error) {
  userId, err := getCurrentUserId()
  if err != nil {
    return nil, err
  }
  return t.GetFollowerIds(userId)
}

func (t TwitterConnection) GetUsersShowById(userId int64) (anaconda.User, error){
  log.Printf("Getting user show %d", userId)
  v := url.Values{}
  result, err := t.api.GetUsersShowById(userId, v)
  log.Printf("Finished Getting user show %d", userId)
  if err != nil {
    log.Printf("Received Error while getting user show");
  } else {
    log.Printf("User name is %s", result.Name)
  }
  return result, err
}

func (t TwitterConnection) GetFollowerIds(userId int64) ([]int64, error) {
  strUserId := strconv.FormatInt(userId, 10)

  v := url.Values{}
  v.Set("user_id", strUserId)
  v.Set("count", "5000")

  c, err := t.api.GetFollowersIds(v)
  if err != nil {
    return nil, err
  }

  return c.Ids, nil
}

func getCurrentUserId() (int64, error) {
  strUserId := os.Getenv("TWITTER_USER_ID")
  return strconv.ParseInt(strUserId, 10, 64)
}
