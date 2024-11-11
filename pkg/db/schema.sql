PRAGMA foreign_keys = on;

-- =======
-- DB meta
-- =======

create table db_version (
    version integer not null
) strict;
insert into db_version values(0);

create table food_types (rowid integer primary key,
    name text not null unique check(length(name) != 0)
) strict;
insert into food_types (name) values
    ('grocery'),
    ('recipe'),
    ('daily log');

create table foods (rowid integer primary key,
    name text not null check(length(name) != 0),

    cals real not null,
    carbs real not null,
    protein real not null,
    fat real not null,

    sugar real not null,
    alcohol real not null default 0,
    water real not null default 0,

    potassium real not null default 0,
    sodium real not null default 0,
    calcium real not null default 0,
    magnesium real not null default 0,
    phosphorus real not null default 0,
    iron real not null default 0,
    zinc real not null default 0,

    mass real not null default 100,
    price real not null default 0,
    density real not null default 1,
    cook_ratio real not null default 1
) strict;


create table units (rowid integer primary key,
    name text not null unique check(length(name) != 0),
    abbreviation text not null unique check(length(abbreviation) != 0)
    -- is_metric integer not null check(is_metric in (0, 1))
) strict;
insert into units(rowid, name, abbreviation) values
    -- Count
    (1, 'count', 'ct'),
    -- Mass
    (2, 'grams', 'g'),
    (3, 'pounds', 'lbs'),
    (4, 'ounces', 'oz'),
    -- Volume
    (5, 'milliliters', 'mL'),
    (6, 'cups', 'cups'),
    (7, 'teaspoons', 'tsp'),
    (8, 'tablespoons', 'tbsp'),
    (9, 'fluid ounces', 'fl-oz');


create table ingredients (rowid integer primary key,
    food_id integer references foods(rowid),
    recipe_id integer references recipes(rowid),

    quantity real not null default 1,
    units integer not null default 0, -- Display purposes only

    in_recipe_id integer references recipes(rowid) on delete cascade not null,
    list_order integer not null,
    is_hidden integer not null default false,
    unique (in_recipe_id, list_order)
    check((food_id is null) + (recipe_id is null) = 1) -- Exactly one should be active
) strict;

create table recipes (rowid integer primary key,
    name text not null check(length(name) != 0),
    blurb text not null,
    instructions text not null,

    computed_food_id integer references foods(rowid) not null
) strict;

create table iterations (rowid integer primary key,
    original_recipe_id integer references recipes(rowid),
    -- original_author integer not null, -- For azimuth integration
    derived_recipe_id integer references recipes(rowid),
    unique(derived_recipe_id)
) strict;

create table daily_logs (rowid integer primary key,
    date integer not null unique,

    computed_food_id integer references foods(rowid) not null
);
