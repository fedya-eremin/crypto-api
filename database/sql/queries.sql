-- name: AddCurrencyToWatchlist :exec
insert into coin (uuid, symbol, watching, interval)
values
($1, $2, true, $3)
on conflict (symbol) do update
set watching = true,
	interval = $3;

-- name: GetNearestPrice :one
select cpl.uuid, cpl.price_usd, cpl.collected_at
from coin_price_log cpl
join coin c on cpl.coin_uuid = c.uuid
where c.symbol = sqlc.arg(symbol) --and cpl.collected_at <= to_timestamp(sqlc.arg(timestamp)::bigint)
order by abs(extract(epoch from (cpl.collected_at - to_timestamp(sqlc.arg(timestamp)::bigint))))
limit 1;

-- name: UnwatchCurrency :exec
update coin
set watching = false
where symbol = $1;

-- name: AddCurrencyPriceLog :exec
insert into coin_price_log (uuid, price_usd, coin_uuid, collected_at)
values
($1, $2, (select uuid from coin where symbol = $3), $4);

-- name: BootstrapWatchingEntries :many
select c.uuid, c.symbol, c.interval from coin c
where c.watching is true;
