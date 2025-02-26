// Code generated by SQLBoiler 4.15.0 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package models

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/friendsofgo/errors"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/volatiletech/sqlboiler/v4/queries/qmhelper"
	"github.com/volatiletech/strmangle"
)

// TeamInfo is an object representing the database table.
type TeamInfo struct {
	ID             int         `boil:"id" json:"id" toml:"id" yaml:"id"`
	Name           string      `boil:"name" json:"name" toml:"name" yaml:"name"`
	Slug           string      `boil:"slug" json:"slug" toml:"slug" yaml:"slug"`
	ShortName      string      `boil:"short_name" json:"short_name" toml:"short_name" yaml:"short_name"`
	Abbr           string      `boil:"abbr" json:"abbr" toml:"abbr" yaml:"abbr"`
	ArenaName      string      `boil:"arena_name" json:"arena_name" toml:"arena_name" yaml:"arena_name"`
	ArenaSize      int         `boil:"arena_size" json:"arena_size" toml:"arena_size" yaml:"arena_size"`
	Color          string      `boil:"color" json:"color" toml:"color" yaml:"color"`
	SecondaryColor string      `boil:"secondary_color" json:"secondary_color" toml:"secondary_color" yaml:"secondary_color"`
	Logo           null.String `boil:"logo" json:"logo,omitempty" toml:"logo" yaml:"logo,omitempty"`

	R *teamInfoR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L teamInfoL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var TeamInfoColumns = struct {
	ID             string
	Name           string
	Slug           string
	ShortName      string
	Abbr           string
	ArenaName      string
	ArenaSize      string
	Color          string
	SecondaryColor string
	Logo           string
}{
	ID:             "id",
	Name:           "name",
	Slug:           "slug",
	ShortName:      "short_name",
	Abbr:           "abbr",
	ArenaName:      "arena_name",
	ArenaSize:      "arena_size",
	Color:          "color",
	SecondaryColor: "secondary_color",
	Logo:           "logo",
}

var TeamInfoTableColumns = struct {
	ID             string
	Name           string
	Slug           string
	ShortName      string
	Abbr           string
	ArenaName      string
	ArenaSize      string
	Color          string
	SecondaryColor string
	Logo           string
}{
	ID:             "team_info.id",
	Name:           "team_info.name",
	Slug:           "team_info.slug",
	ShortName:      "team_info.short_name",
	Abbr:           "team_info.abbr",
	ArenaName:      "team_info.arena_name",
	ArenaSize:      "team_info.arena_size",
	Color:          "team_info.color",
	SecondaryColor: "team_info.secondary_color",
	Logo:           "team_info.logo",
}

// Generated where

var TeamInfoWhere = struct {
	ID             whereHelperint
	Name           whereHelperstring
	Slug           whereHelperstring
	ShortName      whereHelperstring
	Abbr           whereHelperstring
	ArenaName      whereHelperstring
	ArenaSize      whereHelperint
	Color          whereHelperstring
	SecondaryColor whereHelperstring
	Logo           whereHelpernull_String
}{
	ID:             whereHelperint{field: "\"team_info\".\"id\""},
	Name:           whereHelperstring{field: "\"team_info\".\"name\""},
	Slug:           whereHelperstring{field: "\"team_info\".\"slug\""},
	ShortName:      whereHelperstring{field: "\"team_info\".\"short_name\""},
	Abbr:           whereHelperstring{field: "\"team_info\".\"abbr\""},
	ArenaName:      whereHelperstring{field: "\"team_info\".\"arena_name\""},
	ArenaSize:      whereHelperint{field: "\"team_info\".\"arena_size\""},
	Color:          whereHelperstring{field: "\"team_info\".\"color\""},
	SecondaryColor: whereHelperstring{field: "\"team_info\".\"secondary_color\""},
	Logo:           whereHelpernull_String{field: "\"team_info\".\"logo\""},
}

// TeamInfoRels is where relationship names are stored.
var TeamInfoRels = struct {
}{}

// teamInfoR is where relationships are stored.
type teamInfoR struct {
}

// NewStruct creates a new relationship struct
func (*teamInfoR) NewStruct() *teamInfoR {
	return &teamInfoR{}
}

// teamInfoL is where Load methods for each relationship are stored.
type teamInfoL struct{}

var (
	teamInfoAllColumns            = []string{"id", "name", "slug", "short_name", "abbr", "arena_name", "arena_size", "color", "secondary_color", "logo"}
	teamInfoColumnsWithoutDefault = []string{"id", "name", "slug", "short_name", "abbr", "arena_name", "arena_size", "color", "secondary_color"}
	teamInfoColumnsWithDefault    = []string{"logo"}
	teamInfoPrimaryKeyColumns     = []string{"id"}
	teamInfoGeneratedColumns      = []string{}
)

type (
	// TeamInfoSlice is an alias for a slice of pointers to TeamInfo.
	// This should almost always be used instead of []TeamInfo.
	TeamInfoSlice []*TeamInfo
	// TeamInfoHook is the signature for custom TeamInfo hook methods
	TeamInfoHook func(context.Context, boil.ContextExecutor, *TeamInfo) error

	teamInfoQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	teamInfoType                 = reflect.TypeOf(&TeamInfo{})
	teamInfoMapping              = queries.MakeStructMapping(teamInfoType)
	teamInfoPrimaryKeyMapping, _ = queries.BindMapping(teamInfoType, teamInfoMapping, teamInfoPrimaryKeyColumns)
	teamInfoInsertCacheMut       sync.RWMutex
	teamInfoInsertCache          = make(map[string]insertCache)
	teamInfoUpdateCacheMut       sync.RWMutex
	teamInfoUpdateCache          = make(map[string]updateCache)
	teamInfoUpsertCacheMut       sync.RWMutex
	teamInfoUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

var teamInfoAfterSelectHooks []TeamInfoHook

var teamInfoBeforeInsertHooks []TeamInfoHook
var teamInfoAfterInsertHooks []TeamInfoHook

var teamInfoBeforeUpdateHooks []TeamInfoHook
var teamInfoAfterUpdateHooks []TeamInfoHook

var teamInfoBeforeDeleteHooks []TeamInfoHook
var teamInfoAfterDeleteHooks []TeamInfoHook

var teamInfoBeforeUpsertHooks []TeamInfoHook
var teamInfoAfterUpsertHooks []TeamInfoHook

// doAfterSelectHooks executes all "after Select" hooks.
func (o *TeamInfo) doAfterSelectHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range teamInfoAfterSelectHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *TeamInfo) doBeforeInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range teamInfoBeforeInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *TeamInfo) doAfterInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range teamInfoAfterInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *TeamInfo) doBeforeUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range teamInfoBeforeUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *TeamInfo) doAfterUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range teamInfoAfterUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *TeamInfo) doBeforeDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range teamInfoBeforeDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *TeamInfo) doAfterDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range teamInfoAfterDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *TeamInfo) doBeforeUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range teamInfoBeforeUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *TeamInfo) doAfterUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range teamInfoAfterUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddTeamInfoHook registers your hook function for all future operations.
func AddTeamInfoHook(hookPoint boil.HookPoint, teamInfoHook TeamInfoHook) {
	switch hookPoint {
	case boil.AfterSelectHook:
		teamInfoAfterSelectHooks = append(teamInfoAfterSelectHooks, teamInfoHook)
	case boil.BeforeInsertHook:
		teamInfoBeforeInsertHooks = append(teamInfoBeforeInsertHooks, teamInfoHook)
	case boil.AfterInsertHook:
		teamInfoAfterInsertHooks = append(teamInfoAfterInsertHooks, teamInfoHook)
	case boil.BeforeUpdateHook:
		teamInfoBeforeUpdateHooks = append(teamInfoBeforeUpdateHooks, teamInfoHook)
	case boil.AfterUpdateHook:
		teamInfoAfterUpdateHooks = append(teamInfoAfterUpdateHooks, teamInfoHook)
	case boil.BeforeDeleteHook:
		teamInfoBeforeDeleteHooks = append(teamInfoBeforeDeleteHooks, teamInfoHook)
	case boil.AfterDeleteHook:
		teamInfoAfterDeleteHooks = append(teamInfoAfterDeleteHooks, teamInfoHook)
	case boil.BeforeUpsertHook:
		teamInfoBeforeUpsertHooks = append(teamInfoBeforeUpsertHooks, teamInfoHook)
	case boil.AfterUpsertHook:
		teamInfoAfterUpsertHooks = append(teamInfoAfterUpsertHooks, teamInfoHook)
	}
}

// OneG returns a single teamInfo record from the query using the global executor.
func (q teamInfoQuery) OneG(ctx context.Context) (*TeamInfo, error) {
	return q.One(ctx, boil.GetContextDB())
}

// One returns a single teamInfo record from the query.
func (q teamInfoQuery) One(ctx context.Context, exec boil.ContextExecutor) (*TeamInfo, error) {
	o := &TeamInfo{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for team_info")
	}

	if err := o.doAfterSelectHooks(ctx, exec); err != nil {
		return o, err
	}

	return o, nil
}

// AllG returns all TeamInfo records from the query using the global executor.
func (q teamInfoQuery) AllG(ctx context.Context) (TeamInfoSlice, error) {
	return q.All(ctx, boil.GetContextDB())
}

// All returns all TeamInfo records from the query.
func (q teamInfoQuery) All(ctx context.Context, exec boil.ContextExecutor) (TeamInfoSlice, error) {
	var o []*TeamInfo

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to TeamInfo slice")
	}

	if len(teamInfoAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(ctx, exec); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// CountG returns the count of all TeamInfo records in the query using the global executor
func (q teamInfoQuery) CountG(ctx context.Context) (int64, error) {
	return q.Count(ctx, boil.GetContextDB())
}

// Count returns the count of all TeamInfo records in the query.
func (q teamInfoQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count team_info rows")
	}

	return count, nil
}

// ExistsG checks if the row exists in the table using the global executor.
func (q teamInfoQuery) ExistsG(ctx context.Context) (bool, error) {
	return q.Exists(ctx, boil.GetContextDB())
}

// Exists checks if the row exists in the table.
func (q teamInfoQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if team_info exists")
	}

	return count > 0, nil
}

// TeamInfos retrieves all the records using an executor.
func TeamInfos(mods ...qm.QueryMod) teamInfoQuery {
	mods = append(mods, qm.From("\"team_info\""))
	q := NewQuery(mods...)
	if len(queries.GetSelect(q)) == 0 {
		queries.SetSelect(q, []string{"\"team_info\".*"})
	}

	return teamInfoQuery{q}
}

// FindTeamInfoG retrieves a single record by ID.
func FindTeamInfoG(ctx context.Context, iD int, selectCols ...string) (*TeamInfo, error) {
	return FindTeamInfo(ctx, boil.GetContextDB(), iD, selectCols...)
}

// FindTeamInfo retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindTeamInfo(ctx context.Context, exec boil.ContextExecutor, iD int, selectCols ...string) (*TeamInfo, error) {
	teamInfoObj := &TeamInfo{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"team_info\" where \"id\"=$1", sel,
	)

	q := queries.Raw(query, iD)

	err := q.Bind(ctx, exec, teamInfoObj)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from team_info")
	}

	if err = teamInfoObj.doAfterSelectHooks(ctx, exec); err != nil {
		return teamInfoObj, err
	}

	return teamInfoObj, nil
}

// InsertG a single record. See Insert for whitelist behavior description.
func (o *TeamInfo) InsertG(ctx context.Context, columns boil.Columns) error {
	return o.Insert(ctx, boil.GetContextDB(), columns)
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *TeamInfo) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("models: no team_info provided for insertion")
	}

	var err error

	if err := o.doBeforeInsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(teamInfoColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	teamInfoInsertCacheMut.RLock()
	cache, cached := teamInfoInsertCache[key]
	teamInfoInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			teamInfoAllColumns,
			teamInfoColumnsWithDefault,
			teamInfoColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(teamInfoType, teamInfoMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(teamInfoType, teamInfoMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"team_info\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"team_info\" %sDEFAULT VALUES%s"
		}

		var queryOutput, queryReturning string

		if len(cache.retMapping) != 0 {
			queryReturning = fmt.Sprintf(" RETURNING \"%s\"", strings.Join(returnColumns, "\",\""))
		}

		cache.query = fmt.Sprintf(cache.query, queryOutput, queryReturning)
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, vals)
	}

	if len(cache.retMapping) != 0 {
		err = exec.QueryRowContext(ctx, cache.query, vals...).Scan(queries.PtrsFromMapping(value, cache.retMapping)...)
	} else {
		_, err = exec.ExecContext(ctx, cache.query, vals...)
	}

	if err != nil {
		return errors.Wrap(err, "models: unable to insert into team_info")
	}

	if !cached {
		teamInfoInsertCacheMut.Lock()
		teamInfoInsertCache[key] = cache
		teamInfoInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(ctx, exec)
}

// UpdateG a single TeamInfo record using the global executor.
// See Update for more documentation.
func (o *TeamInfo) UpdateG(ctx context.Context, columns boil.Columns) (int64, error) {
	return o.Update(ctx, boil.GetContextDB(), columns)
}

// Update uses an executor to update the TeamInfo.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *TeamInfo) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	var err error
	if err = o.doBeforeUpdateHooks(ctx, exec); err != nil {
		return 0, err
	}
	key := makeCacheKey(columns, nil)
	teamInfoUpdateCacheMut.RLock()
	cache, cached := teamInfoUpdateCache[key]
	teamInfoUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			teamInfoAllColumns,
			teamInfoPrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("models: unable to update team_info, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"team_info\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, teamInfoPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(teamInfoType, teamInfoMapping, append(wl, teamInfoPrimaryKeyColumns...))
		if err != nil {
			return 0, err
		}
	}

	values := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), cache.valueMapping)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, values)
	}
	var result sql.Result
	result, err = exec.ExecContext(ctx, cache.query, values...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update team_info row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by update for team_info")
	}

	if !cached {
		teamInfoUpdateCacheMut.Lock()
		teamInfoUpdateCache[key] = cache
		teamInfoUpdateCacheMut.Unlock()
	}

	return rowsAff, o.doAfterUpdateHooks(ctx, exec)
}

// UpdateAllG updates all rows with the specified column values.
func (q teamInfoQuery) UpdateAllG(ctx context.Context, cols M) (int64, error) {
	return q.UpdateAll(ctx, boil.GetContextDB(), cols)
}

// UpdateAll updates all rows with the specified column values.
func (q teamInfoQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all for team_info")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected for team_info")
	}

	return rowsAff, nil
}

// UpdateAllG updates all rows with the specified column values.
func (o TeamInfoSlice) UpdateAllG(ctx context.Context, cols M) (int64, error) {
	return o.UpdateAll(ctx, boil.GetContextDB(), cols)
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o TeamInfoSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	ln := int64(len(o))
	if ln == 0 {
		return 0, nil
	}

	if len(cols) == 0 {
		return 0, errors.New("models: update all requires at least one column argument")
	}

	colNames := make([]string, len(cols))
	args := make([]interface{}, len(cols))

	i := 0
	for name, value := range cols {
		colNames[i] = name
		args[i] = value
		i++
	}

	// Append all of the primary key values for each column
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), teamInfoPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"team_info\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, teamInfoPrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all in teamInfo slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected all in update all teamInfo")
	}
	return rowsAff, nil
}

// UpsertG attempts an insert, and does an update or ignore on conflict.
func (o *TeamInfo) UpsertG(ctx context.Context, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	return o.Upsert(ctx, boil.GetContextDB(), updateOnConflict, conflictColumns, updateColumns, insertColumns)
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *TeamInfo) Upsert(ctx context.Context, exec boil.ContextExecutor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("models: no team_info provided for upsert")
	}

	if err := o.doBeforeUpsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(teamInfoColumnsWithDefault, o)

	// Build cache key in-line uglily - mysql vs psql problems
	buf := strmangle.GetBuffer()
	if updateOnConflict {
		buf.WriteByte('t')
	} else {
		buf.WriteByte('f')
	}
	buf.WriteByte('.')
	for _, c := range conflictColumns {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	buf.WriteString(strconv.Itoa(updateColumns.Kind))
	for _, c := range updateColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	buf.WriteString(strconv.Itoa(insertColumns.Kind))
	for _, c := range insertColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	for _, c := range nzDefaults {
		buf.WriteString(c)
	}
	key := buf.String()
	strmangle.PutBuffer(buf)

	teamInfoUpsertCacheMut.RLock()
	cache, cached := teamInfoUpsertCache[key]
	teamInfoUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			teamInfoAllColumns,
			teamInfoColumnsWithDefault,
			teamInfoColumnsWithoutDefault,
			nzDefaults,
		)

		update := updateColumns.UpdateColumnSet(
			teamInfoAllColumns,
			teamInfoPrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("models: unable to upsert team_info, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(teamInfoPrimaryKeyColumns))
			copy(conflict, teamInfoPrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"team_info\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(teamInfoType, teamInfoMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(teamInfoType, teamInfoMapping, ret)
			if err != nil {
				return err
			}
		}
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)
	var returns []interface{}
	if len(cache.retMapping) != 0 {
		returns = queries.PtrsFromMapping(value, cache.retMapping)
	}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, vals)
	}
	if len(cache.retMapping) != 0 {
		err = exec.QueryRowContext(ctx, cache.query, vals...).Scan(returns...)
		if errors.Is(err, sql.ErrNoRows) {
			err = nil // Postgres doesn't return anything when there's no update
		}
	} else {
		_, err = exec.ExecContext(ctx, cache.query, vals...)
	}
	if err != nil {
		return errors.Wrap(err, "models: unable to upsert team_info")
	}

	if !cached {
		teamInfoUpsertCacheMut.Lock()
		teamInfoUpsertCache[key] = cache
		teamInfoUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(ctx, exec)
}

// DeleteG deletes a single TeamInfo record.
// DeleteG will match against the primary key column to find the record to delete.
func (o *TeamInfo) DeleteG(ctx context.Context) (int64, error) {
	return o.Delete(ctx, boil.GetContextDB())
}

// Delete deletes a single TeamInfo record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *TeamInfo) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("models: no TeamInfo provided for delete")
	}

	if err := o.doBeforeDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), teamInfoPrimaryKeyMapping)
	sql := "DELETE FROM \"team_info\" WHERE \"id\"=$1"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete from team_info")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by delete for team_info")
	}

	if err := o.doAfterDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	return rowsAff, nil
}

func (q teamInfoQuery) DeleteAllG(ctx context.Context) (int64, error) {
	return q.DeleteAll(ctx, boil.GetContextDB())
}

// DeleteAll deletes all matching rows.
func (q teamInfoQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("models: no teamInfoQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from team_info")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for team_info")
	}

	return rowsAff, nil
}

// DeleteAllG deletes all rows in the slice.
func (o TeamInfoSlice) DeleteAllG(ctx context.Context) (int64, error) {
	return o.DeleteAll(ctx, boil.GetContextDB())
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o TeamInfoSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	if len(teamInfoBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), teamInfoPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"team_info\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, teamInfoPrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from teamInfo slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for team_info")
	}

	if len(teamInfoAfterDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	return rowsAff, nil
}

// ReloadG refetches the object from the database using the primary keys.
func (o *TeamInfo) ReloadG(ctx context.Context) error {
	if o == nil {
		return errors.New("models: no TeamInfo provided for reload")
	}

	return o.Reload(ctx, boil.GetContextDB())
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *TeamInfo) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindTeamInfo(ctx, exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAllG refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *TeamInfoSlice) ReloadAllG(ctx context.Context) error {
	if o == nil {
		return errors.New("models: empty TeamInfoSlice provided for reload all")
	}

	return o.ReloadAll(ctx, boil.GetContextDB())
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *TeamInfoSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := TeamInfoSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), teamInfoPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"team_info\".* FROM \"team_info\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, teamInfoPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in TeamInfoSlice")
	}

	*o = slice

	return nil
}

// TeamInfoExistsG checks if the TeamInfo row exists.
func TeamInfoExistsG(ctx context.Context, iD int) (bool, error) {
	return TeamInfoExists(ctx, boil.GetContextDB(), iD)
}

// TeamInfoExists checks if the TeamInfo row exists.
func TeamInfoExists(ctx context.Context, exec boil.ContextExecutor, iD int) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"team_info\" where \"id\"=$1 limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, iD)
	}
	row := exec.QueryRowContext(ctx, sql, iD)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if team_info exists")
	}

	return exists, nil
}

// Exists checks if the TeamInfo row exists.
func (o *TeamInfo) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	return TeamInfoExists(ctx, exec, o.ID)
}
