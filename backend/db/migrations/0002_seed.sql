-- +goose Up

INSERT INTO entity_types (key, name_en, name_de, icon, color, is_system) VALUES
  ('character', 'Character', 'Charakter',  'user',           '#a8442e', 1),
  ('location',  'Location',  'Ort',        'map-pin',        '#173842', 1),
  ('faction',   'Faction',   'Fraktion',   'flag',           '#7c5e30', 1),
  ('item',      'Item',      'Gegenstand', 'package',        '#b8935a', 1),
  ('deity',     'Deity',     'Gottheit',   'star',           '#d4b06a', 1),
  ('race',      'Race',      'Volk',       'users',          '#5c4a3a', 1),
  ('creature',  'Creature',  'Kreatur',    'feather',        '#7b2b1a', 1),
  ('religion',  'Religion',  'Religion',   'book',           '#173842', 1),
  ('language',  'Language',  'Sprache',    'message-circle', '#7c5e30', 1),
  ('culture',   'Culture',   'Kultur',     'globe',          '#a8442e', 1),
  ('concept',   'Concept',   'Konzept',    'sparkles',       '#b8935a', 1);

INSERT INTO field_definitions (entity_type_id, key, label_en, label_de, data_type, config, sort_order, is_required) VALUES
  ((SELECT id FROM entity_types WHERE key='character'), 'race',        'Race',        'Volk',           'entity_ref', '{"target_type":"race"}',    1, 0),
  ((SELECT id FROM entity_types WHERE key='character'), 'status',      'Status',      'Status',         'select',     '{"options":["alive","dead","unknown","missing"]}', 2, 0),
  ((SELECT id FROM entity_types WHERE key='character'), 'affiliation', 'Affiliation', 'Zugehörigkeit',  'entity_ref', '{"target_type":"faction"}', 3, 0),
  ((SELECT id FROM entity_types WHERE key='character'), 'born',        'Born',        'Geboren',        'date',       NULL, 4, 0),
  ((SELECT id FROM entity_types WHERE key='character'), 'died',        'Died',        'Gestorben',      'date',       NULL, 5, 0),

  ((SELECT id FROM entity_types WHERE key='location'),  'population',  'Population',  'Bevölkerung',    'number',     NULL, 1, 0),
  ((SELECT id FROM entity_types WHERE key='location'),  'climate',     'Climate',     'Klima',          'text',       NULL, 2, 0),
  ((SELECT id FROM entity_types WHERE key='location'),  'government',  'Government',  'Regierung',      'text',       NULL, 3, 0),

  ((SELECT id FROM entity_types WHERE key='faction'),   'headquarters','Headquarters','Hauptsitz',      'entity_ref', '{"target_type":"location"}', 1, 0),
  ((SELECT id FROM entity_types WHERE key='faction'),   'founded',     'Founded',     'Gegründet',      'date',       NULL, 2, 0),

  ((SELECT id FROM entity_types WHERE key='item'),      'rarity',      'Rarity',      'Seltenheit',     'select',     '{"options":["common","uncommon","rare","mythic"]}', 1, 0),
  ((SELECT id FROM entity_types WHERE key='item'),      'material',    'Material',    'Material',       'text',       NULL, 2, 0),

  ((SELECT id FROM entity_types WHERE key='deity'),     'domain',      'Domain',      'Domäne',         'text',       NULL, 1, 0),
  ((SELECT id FROM entity_types WHERE key='deity'),     'symbol',      'Symbol',      'Symbol',         'text',       NULL, 2, 0),

  ((SELECT id FROM entity_types WHERE key='race'),      'origin',      'Origin',      'Ursprung',       'text',       NULL, 1, 0),
  ((SELECT id FROM entity_types WHERE key='race'),      'lifespan',    'Lifespan',    'Lebenserwartung','text',       NULL, 2, 0);

INSERT INTO relationship_types (key, label_en, label_de, inverse_label_en, inverse_label_de, is_symmetric, category) VALUES
  ('ally_of',    'ally of',    'verbündet mit',  'ally of',     'verbündet mit',    1, 'generic'),
  ('enemy_of',   'enemy of',   'verfeindet mit', 'enemy of',    'verfeindet mit',   1, 'generic'),
  ('member_of',  'member of',  'Mitglied von',   'has member',  'hat Mitglied',     0, 'generic'),
  ('owns',       'owns',       'besitzt',        'owned by',    'besessen von',     0, 'generic'),
  ('located_in', 'located in', 'liegt in',       'contains',    'enthält',          0, 'generic'),
  ('created_by', 'created by', 'erschaffen von', 'creator of',  'Schöpfer von',     0, 'generic'),
  ('founded_by', 'founded by', 'gegründet von',  'founder of',  'Gründer von',      0, 'generic'),
  ('parent_of',  'parent of',  'Elternteil von', 'child of',    'Kind von',         0, 'genealogy'),
  ('spouse_of',  'spouse of',  'verheiratet mit','spouse of',   'verheiratet mit',  1, 'genealogy'),
  ('sibling_of', 'sibling of', 'Geschwister von','sibling of',  'Geschwister von',  1, 'genealogy');

INSERT INTO calendars (name, epoch_name, current_year, current_month_index, current_day) VALUES
  ('Luminox Reckoning', 'Luminox', 1187, 4, 14);

INSERT INTO calendar_months (calendar_id, sort_order, name, days) VALUES
  ((SELECT id FROM calendars WHERE name='Luminox Reckoning'), 1, 'Frostbell',  36),
  ((SELECT id FROM calendars WHERE name='Luminox Reckoning'), 2, 'Miravel',    36),
  ((SELECT id FROM calendars WHERE name='Luminox Reckoning'), 3, 'Ashtide',    36),
  ((SELECT id FROM calendars WHERE name='Luminox Reckoning'), 4, 'Veilmoon',   36),
  ((SELECT id FROM calendars WHERE name='Luminox Reckoning'), 5, 'Sundering',  36),
  ((SELECT id FROM calendars WHERE name='Luminox Reckoning'), 6, 'Deepgloam',  36),
  ((SELECT id FROM calendars WHERE name='Luminox Reckoning'), 7, 'Wane',       36),
  ((SELECT id FROM calendars WHERE name='Luminox Reckoning'), 8, 'Luminox',    36),
  ((SELECT id FROM calendars WHERE name='Luminox Reckoning'), 9, 'Harrow',     36);

INSERT INTO calendar_weekdays (calendar_id, sort_order, name) VALUES
  ((SELECT id FROM calendars WHERE name='Luminox Reckoning'), 1, 'Cinderday'),
  ((SELECT id FROM calendars WHERE name='Luminox Reckoning'), 2, 'Mireday'),
  ((SELECT id FROM calendars WHERE name='Luminox Reckoning'), 3, 'Veil'),
  ((SELECT id FROM calendars WHERE name='Luminox Reckoning'), 4, 'Songday'),
  ((SELECT id FROM calendars WHERE name='Luminox Reckoning'), 5, 'Hollow'),
  ((SELECT id FROM calendars WHERE name='Luminox Reckoning'), 6, 'Stillnight');

INSERT INTO calendar_eras (calendar_id, name, start_year, suffix) VALUES
  ((SELECT id FROM calendars WHERE name='Luminox Reckoning'), 'Before the Sundering', -2000, 'BS'),
  ((SELECT id FROM calendars WHERE name='Luminox Reckoning'), 'Luminox Reckoning',    0,    'LR');

INSERT INTO calendar_moons (calendar_id, name, cycle_days, offset_days) VALUES
  ((SELECT id FROM calendars WHERE name='Luminox Reckoning'), 'Vael',   28, 0),
  ((SELECT id FROM calendars WHERE name='Luminox Reckoning'), 'Miraël', 44, 12),
  ((SELECT id FROM calendars WHERE name='Luminox Reckoning'), 'Korren', 61, 0);

INSERT INTO entities (entity_type_id, title, slug, summary, body, body_text, visibility, created_at, updated_at) VALUES
  ((SELECT id FROM entity_types WHERE key='location'),  'Halvenmoor', 'halvenmoor',
   'A drowned realm in the north — the homeland of the Halvenkin and birthplace of Cassara.',
   '{"type":"doc","content":[{"type":"paragraph","content":[{"type":"text","text":"Halvenmoor lies where the great rivers meet the salt — a country half-water, half-land. Its court fell during the Sundering; its songs survive."}]}]}',
   'Halvenmoor lies where the great rivers meet the salt — a country half-water, half-land. Its court fell during the Sundering; its songs survive.',
   'player', strftime('%Y-%m-%dT%H:%M:%fZ','now'), strftime('%Y-%m-%dT%H:%M:%fZ','now')),

  ((SELECT id FROM entity_types WHERE key='location'),  'Varkuun', 'varkuun',
   'The seat of the Lazarine Court, perched between two basalt cliffs.',
   '{"type":"doc","content":[{"type":"paragraph","content":[{"type":"text","text":"Varkuun was built atop the bones of a god whose name is no longer pronounced. Its towers are roped with iron and silk."}]}]}',
   'Varkuun was built atop the bones of a god whose name is no longer pronounced. Its towers are roped with iron and silk.',
   'player', strftime('%Y-%m-%dT%H:%M:%fZ','now'), strftime('%Y-%m-%dT%H:%M:%fZ','now')),

  ((SELECT id FROM entity_types WHERE key='race'),      'Halvenkin', 'halvenkin',
   'The marsh-folk of Halvenmoor; pale, long-fingered, and slow to laugh.',
   '{"type":"doc","content":[{"type":"paragraph","content":[{"type":"text","text":"The Halvenkin trace their lineage to the drowning of Halvenmoor — a people remade by tide and salt."}]}]}',
   'The Halvenkin trace their lineage to the drowning of Halvenmoor — a people remade by tide and salt.',
   'player', strftime('%Y-%m-%dT%H:%M:%fZ','now'), strftime('%Y-%m-%dT%H:%M:%fZ','now')),

  ((SELECT id FROM entity_types WHERE key='deity'),     'Vael, the Counter', 'vael-the-counter',
   'A small god of the drowned, who counts each soul that crosses the shoal.',
   '{"type":"doc","content":[{"type":"paragraph","content":[{"type":"text","text":"Vael is no longer worshipped openly. The Lazarine Court keeps a single chapel for her."}]}]}',
   'Vael is no longer worshipped openly. The Lazarine Court keeps a single chapel for her.',
   'player', strftime('%Y-%m-%dT%H:%M:%fZ','now'), strftime('%Y-%m-%dT%H:%M:%fZ','now')),

  ((SELECT id FROM entity_types WHERE key='faction'),   'Lazarine Court', 'lazarine-court',
   'The court of Varkuun — what remains of it after the Sundering.',
   '{"type":"doc","content":[{"type":"paragraph","content":[{"type":"text","text":"Once a sovereignty, now a fellowship — twelve houses, three of which exist only in name."}]}]}',
   'Once a sovereignty, now a fellowship — twelve houses, three of which exist only in name.',
   'player', strftime('%Y-%m-%dT%H:%M:%fZ','now'), strftime('%Y-%m-%dT%H:%M:%fZ','now')),

  ((SELECT id FROM entity_types WHERE key='character'), 'Veiled Cassara', 'veiled-cassara',
   'Cartographer-priestess of the Drowned Wheel, last living interpreter of the Sundering charts.',
   '{"type":"doc","content":[{"type":"paragraph","content":[{"type":"text","text":"She came north out of "},{"type":"wikilink","attrs":{"slug":"halvenmoor","label":"Halvenmoor"}},{"type":"text","text":" with a single chart folded into the lining of her cloak."}]},{"type":"paragraph","content":[{"type":"text","text":"The court of "},{"type":"wikilink","attrs":{"slug":"varkuun","label":"Varkuun"}},{"type":"text","text":" receives her each Ashtide. She speaks little."}]}]}',
   'She came north out of Halvenmoor with a single chart folded into the lining of her cloak. The court of Varkuun receives her each Ashtide.',
   'player', strftime('%Y-%m-%dT%H:%M:%fZ','now'), strftime('%Y-%m-%dT%H:%M:%fZ','now'));

INSERT INTO tags (name) VALUES
  ('cartographer'), ('lazarine'), ('drowned-wheel'), ('halvenkin'), ('sundering'), ('priesthood');

INSERT INTO entity_tags (entity_id, tag_id)
  SELECT e.id, t.id FROM entities e JOIN tags t
   ON t.name IN ('cartographer', 'lazarine', 'drowned-wheel', 'priesthood')
  WHERE e.slug = 'veiled-cassara';

INSERT INTO entity_tags (entity_id, tag_id)
  SELECT e.id, t.id FROM entities e JOIN tags t
   ON t.name IN ('halvenkin', 'sundering')
  WHERE e.slug = 'halvenmoor';

INSERT INTO entity_tags (entity_id, tag_id)
  SELECT e.id, t.id FROM entities e JOIN tags t
   ON t.name IN ('halvenkin', 'sundering')
  WHERE e.slug = 'halvenkin';

INSERT INTO entity_tags (entity_id, tag_id)
  SELECT e.id, t.id FROM entities e JOIN tags t
   ON t.name = 'lazarine'
  WHERE e.slug = 'lazarine-court';

INSERT INTO links (source_entity_id, target_entity_id, anchor)
  SELECT a.id, b.id, '' FROM entities a, entities b
  WHERE a.slug = 'veiled-cassara' AND b.slug = 'halvenmoor';

INSERT INTO links (source_entity_id, target_entity_id, anchor)
  SELECT a.id, b.id, '' FROM entities a, entities b
  WHERE a.slug = 'veiled-cassara' AND b.slug = 'varkuun';

-- +goose Down

DELETE FROM links;
DELETE FROM entity_tags;
DELETE FROM tags;
DELETE FROM entities;
DELETE FROM calendar_moons;
DELETE FROM calendar_eras;
DELETE FROM calendar_weekdays;
DELETE FROM calendar_months;
DELETE FROM calendars;
DELETE FROM relationship_types;
DELETE FROM field_definitions;
DELETE FROM entity_types;
