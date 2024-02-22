-- +goose Up
-- +goose StatementBegin
UPDATE teams SET display_name = 'Timberwolves' WHERE rsc_id = 1;
UPDATE teams SET display_name = 'Pacers' WHERE rsc_id = 2;
UPDATE teams SET display_name = 'Jazz' WHERE rsc_id = 3;
UPDATE teams SET display_name = 'Magic' WHERE rsc_id = 4;
UPDATE teams SET display_name = 'Hawks' WHERE rsc_id = 5;
UPDATE teams SET display_name = 'Celtics' WHERE rsc_id = 6;
UPDATE teams SET display_name = 'Cavaliers' WHERE rsc_id = 7;
UPDATE teams SET display_name = 'Knicks' WHERE rsc_id = 8;
UPDATE teams SET display_name = 'Pelicans' WHERE rsc_id = 9;
UPDATE teams SET display_name = 'Trail Blazers' WHERE rsc_id = 10;
UPDATE teams SET display_name = 'Grizzlies' WHERE rsc_id = 11;
UPDATE teams SET display_name = 'Lakers' WHERE rsc_id = 12;
UPDATE teams SET display_name = 'Thunder' WHERE rsc_id = 13;
UPDATE teams SET display_name = 'Mavericks' WHERE rsc_id = 14;
UPDATE teams SET display_name = 'Rockets' WHERE rsc_id = 15;
UPDATE teams SET display_name = 'Nuggets' WHERE rsc_id = 16;
UPDATE teams SET display_name = '76ers' WHERE rsc_id = 17;
UPDATE teams SET display_name = 'Nets' WHERE rsc_id = 18;
UPDATE teams SET display_name = 'Kings' WHERE rsc_id = 19;
UPDATE teams SET display_name = 'Heat' WHERE rsc_id = 20;
UPDATE teams SET display_name = 'Warriors' WHERE rsc_id = 21;
UPDATE teams SET display_name = 'Bulls' WHERE rsc_id = 22;
UPDATE teams SET display_name = 'Clippers' WHERE rsc_id = 23;
UPDATE teams SET display_name = 'Suns' WHERE rsc_id = 24;
UPDATE teams SET display_name = 'Bucks' WHERE rsc_id = 25;
UPDATE teams SET display_name = 'Pistons' WHERE rsc_id = 26;
UPDATE teams SET display_name = 'Hornets' WHERE rsc_id = 27;
UPDATE teams SET display_name = 'Spurs' WHERE rsc_id = 28;
UPDATE teams SET display_name = 'Wizards' WHERE rsc_id = 29;
UPDATE teams SET display_name = 'Raptors' WHERE rsc_id = 30;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
UPDATE teams SET display_name = NULL
-- +goose StatementEnd
