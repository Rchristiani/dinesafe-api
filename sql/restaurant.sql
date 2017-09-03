CREATE TABLE restaurants(
	id                        serial,
	establishmentID           int,
	inspectionID              int,
	establishmentName         varchar,
	establishmentType         varchar,
	establishmentAddress      varchar,
	establishmentStatus       varchar,
	minimumInspectionsPerYear int,
	infractionDetails         varchar,
	inspectionDate            varchar,
	severity                  varchar,
	action                    varchar,
	courtOutcome              varchar,
	amountFinded              varchar
)
