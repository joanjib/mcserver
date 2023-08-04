-- Mcserver
-- Copyright (C) 2023  JUAN JOSÉ IGLESIAS BLANCH

-- This program is free software: you can redistribute it and/or modify
-- it under the terms of the GNU Affero General Public License as published by
-- the Free Software Foundation, either version 3 of the License, or
-- (at your option) any later version.

-- This program is distributed in the hope that it will be useful,
-- but WITHOUT ANY WARRANTY; without even the implied warranty of
-- MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
-- GNU Affero General Public License for more details.

-- You should have received a copy of the GNU Affero General Public License
-- along with this program.  If not, see <http://www.gnu.org/licenses/>.
-- sqlite
-- tables for the tasks manager

PRAGMA foreign_keys = ON;
-- places: every position staff can reach it
CREATE TABLE IF NOT EXISTS location (
    id      integer primary key,
    name    text    not null,
    x       integer not null,
    y       integer not null

);
-- distances: pre-calculated distances calculation between each place
-- origin: identifier of the origin place
-- destin: identifier of the destination place.
CREATE TABLE if not EXISTS distances (
    origin_id   integer not null,
    destin_id   integer not null,
    distance    integer not null, 
    FOREIGN KEY(origin_id)  REFERENCES location(id),   
    FOREIGN KEY(destin_id)  REFERENCES location(id),   
    PRIMARY KEY(origin_id,destin_id)
) WITHOUT ROWID;

CREATE TABLE if not EXISTS scheduling (
    id          integer primary key,
    stating_at  text    not null,
    ending_at   text    not null,
    periodicity integer,    -- 0: not specified, 1: all days, 2: weekly, 3: monthly, 4: yearly, 5: punctual
    planning    text        -- text specifying depending on if is all days, weeksly, yearly, etc...
    -- example: Weekly  -> MWF : Monday, Wendnesday, and friday
    --          Monthly -> 1: means the first, -1 : the last, and so on
    --          Yearly  -> 3,-1: the March, the last day of March
);

CREATE TABLE if not EXISTS task (

    id          integer primary key,
    name        text    not null,
    description text    not null,
    priority    integer not NULL,       -- default priority when planned. Subtasks inherit this priority
    sched_id    integer not null,       -- reference to the schedule of the task 
    FOREIGN KEY(sched_id)  REFERENCES scheduling(id)
);

CREATE TABLE if not EXISTS category (
    id          integer primary key,
    name        text    not null,
    description text    not null
);
CREATE TABLE if not EXISTS staff (
    id          integer primary key,
    username    text    not null unique,
    name        text    not null,
    surname1    text    not null,
    surname2    text,
    nfc_device  text,
    nfc_creden  text    not null,    
    state       integer default 0   -- 0 offline not elegible, 1 online
   
);
CREATE TABLE if not EXISTS patients (
    id          integer primary key,
    name        text    not null,
    surname1    text    not null,
    surname2    text,
    identif     text,                   -- document with unique identification, like dni   
    room_id     integer not null,
    bed         integer not NULL,       -- bed number into the room
    cur_loc_id  integer,                -- current location for restricted in movement patients
    FOREIGN KEY(room_id)    REFERENCES location(id),
    FOREIGN KEY(cur_loc_id) REFERENCES location(id)
);

CREATE TABLE if not EXISTS category_staff (
    cat_id      integer not null,
    staff_id    integer not null,
    FOREIGN KEY(cat_id)  REFERENCES category(id),
    FOREIGN KEY(staff_id)  REFERENCES staff(id),
    PRIMARY KEY(cat_id,staff_id)
) WITHOUT ROWID;

CREATE TABLE if not EXISTS category_patient (
    name        text    not null,
    description text    not null,
    cat_id      integer not null,
    patient_id  integer not null,
    FOREIGN KEY(cat_id)  REFERENCES category(id),
    FOREIGN KEY(patient_id)  REFERENCES patient(id),
    PRIMARY KEY(cat_id,patient_id)
) WITHOUT ROWID;

CREATE TABLE if not EXISTS category_location (
    name        text    not null,
    description text    not null,
    cat_id      integer not null,
    location_id integer not null,
    FOREIGN KEY(cat_id)  REFERENCES category(id),
    FOREIGN KEY(location_id)  REFERENCES location(id),
    PRIMARY KEY(cat_id,location_id)
) WITHOUT ROWID;


CREATE TABLE if not EXISTS subtask (

    id          integer primary key,
    name        text    not null,
    description text    not null,
    voice_comm  text    not null,
    num_seq     integer not null,       -- all tasks in same numeric sequence can be executed in parallel
    task_id     integer not null,
    is_optional integer not null,       -- 0 is mandatory, 1 is optional
    start_time  text    not null,       -- this task can be done starting at this time
    end_time    text    not null,       -- this task has to be completed at this time in worst case.
    for_who_id  integer not null,       -- for who has to be done the task
    to_who_id   integer,                -- to which category the task has to be done    
    location_id integer,
    FOREIGN KEY(task_id)        REFERENCES task(id),
    FOREIGN KEY(to_who_id)      REFERENCES category(id),
    FOREIGN KEY(location_id)    REFERENCES catetory(id)
    FOREIGN KEY(for_who_id)     REFERENCES category(id)

);

CREATE TABLE if not EXISTS planned_task (
    id                      integer primary key,
    starting_date           text    not null,
    ending_date             text    not null,
    actual_starting_date    text    not null,
    actual_ending_date      text    not null,
    task_id                 integer not null,
    state                   integer not NULL,       -- 0: pending, 1:started, 2:done, 3: incidence
    priority                integer not null,
    FOREIGN KEY(task_id)    REFERENCES task(id)
);

CREATE TABLE if not EXISTS planned_subtask (
    id                      integer primary key,
    starting_date           text    not null,
    ending_date             text    not null,
    actual_starting_date    text    not null,
    actual_ending_date      text    not null,
    num_seq                 integer not null,       -- all tasks in same numeric sequence can be executed in parallel
    planned_task_id         integer not null,
    subtask_id              integer not null,
    is_optional             integer not null,       -- 0 is mandatory, 1 is optional
    state                   integer not NULL,       -- 0: pending, 1:started, 2:done, 3: incidence
    priority                integer not null,
    for_who_id              integer not null,       -- for who has to be done the task
    to_who_id               integer,                -- to which category the task has to be done    
    location_id             integer,
    FOREIGN KEY(to_who_id)          REFERENCES patients(id),
    FOREIGN KEY(location_id)        REFERENCES location(id)
    FOREIGN KEY(for_who_id)         REFERENCES staff(id)
    FOREIGN KEY(planned_task_id)    REFERENCES planned_task(id)
    FOREIGN KEY(subtask_id)         REFERENCES subtask(id)

);
