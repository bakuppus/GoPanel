package botcontext

import (
	"fmt"
	"github.com/TicketsBot/GoPanel/config"
	dbclient "github.com/TicketsBot/GoPanel/database"
	"github.com/TicketsBot/GoPanel/redis"
	"github.com/rxdn/gdl/rest/ratelimit"
)

func ContextForGuild(guildId uint64) (ctx BotContext, err error) {
	whitelabelBotId, isWhitelabel, err := dbclient.Client.WhitelabelGuilds.GetBotByGuild(guildId)
	if err != nil {
		return
	}

	if isWhitelabel {
		res, err := dbclient.Client.Whitelabel.GetByBotId(whitelabelBotId)
		if err != nil {
			return ctx, err
		}

		rateLimiter := ratelimit.NewRateLimiter(ratelimit.NewRedisStore(redis.Client.Client, fmt.Sprintf("ratelimiter:%d", whitelabelBotId)), 1)

		return BotContext{
			BotId:       res.BotId,
			Token:       res.Token,
			RateLimiter: rateLimiter,
		}, nil
	} else {
		return PublicContext(), nil
	}
}

func PublicContext() BotContext {
	return BotContext{
		BotId:       config.Conf.Bot.Id,
		Token:       config.Conf.Bot.Token,
		RateLimiter: ratelimit.NewRateLimiter(ratelimit.NewRedisStore(redis.Client.Client, "ratelimiter:public"), 1),
	}
}
