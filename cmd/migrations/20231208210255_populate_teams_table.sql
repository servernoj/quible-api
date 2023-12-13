-- +goose Up
-- +goose StatementBegin
INSERT INTO teams ("name",arena,color,rsc_id) VALUES
	 ('Atlanta Hawks','State Farm Arena','#C8102E',5),
	 ('Boston Celtics','TD Garden','#007A33',6),
	 ('Brooklyn Nets','Barclays Center','#000000',18),
	 ('Charlotte Hornets','Spectrum Center','#1D1160',27),
	 ('Chicago Bulls','United Center','#CE1141',22),
	 ('Cleveland Cavaliers','Rocket Mortgage Field House','#860038',7),
	 ('Dallas Mavericks','American Airlines Center','#00538C',14),
	 ('Denver Nuggets','Ball Arena','#0E2240',16),
	 ('Detroit Pistons','Little Caesars Arena','#C8102E',26),
	 ('Golden State Warriors','Chase Center','#1D428A',21),
	 ('Houston Rockets','Toyota Center','#CE1141',15),
	 ('Indiana Pacers','Bankers Life Fieldhouse','#002D62',2),
	 ('Los Angeles Clippers','Crypto.com Arena','#C8102E',23),
	 ('Los Angeles Lakers','Crypto.com Arena','#552583',12),
	 ('Memphis Grizzlies','FedExForum','#5D76A9',11),
	 ('Miami Heat','FTX Arena','#98002E',20),
	 ('Milwaukee Bucks','Fiserv Forum','#00471B',25),
	 ('Minnesota Timberwolves','Target Center','#0C2340',1),
	 ('New Orleans Pelicans','Smoothie King Center','#0C2340',9),
	 ('New York Knicks','Madison Square Garden','#006BB6',8),
	 ('Oklahoma City Thunder','Paycom Center','#007AC1',13),
	 ('Orlando Magic','Amway Center','#0077C0',4),
	 ('Philadelphia 76ers','Wells Fargo Center','#006BB6',17),
	 ('Phoenix Suns','Footprint Center','#1D1160',24),
	 ('Portland Trail Blazers','Moda Center','#E03A3E',10),
	 ('Sacramento Kings','Golden 1 Center','#5A2D81',19),
	 ('San Antonio Spurs','Frost Bank Center','#C4CED4',28),
	 ('Toronto Raptors','Scotiabank Arena','#CE1141',30),
	 ('Utah Jazz','Vivint Arena','#002B5C',3),
	 ('Washington Wizards','Capital One Arena','#002B5C',29);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
delete from teams;
-- +goose StatementEnd
