create table discord_user (
    id                         serial PRIMARY KEY, 
    discord_user_id            varchar(64) not null,
    discord_username           varchar(128),
    discord_email              varchar(128),
    nftkeyme_id                varchar(128),
    nftkeyme_email             varchar(128),
    nftkeyme_access_token      varchar(128),
    nftkeyme_refresh_token     varchar(128),
    num_assets                 integer,
    UNIQUE(discord_user_id)
);