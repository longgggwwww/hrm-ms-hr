// Code generated by ent, DO NOT EDIT.

package organization

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/longgggwwww/hrm-ms-hr/ent/predicate"
)

// ID filters vertices based on their ID field.
func ID(id int) predicate.Organization {
	return predicate.Organization(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id int) predicate.Organization {
	return predicate.Organization(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id int) predicate.Organization {
	return predicate.Organization(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...int) predicate.Organization {
	return predicate.Organization(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...int) predicate.Organization {
	return predicate.Organization(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id int) predicate.Organization {
	return predicate.Organization(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id int) predicate.Organization {
	return predicate.Organization(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id int) predicate.Organization {
	return predicate.Organization(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id int) predicate.Organization {
	return predicate.Organization(sql.FieldLTE(FieldID, id))
}

// Name applies equality check predicate on the "name" field. It's identical to NameEQ.
func Name(v string) predicate.Organization {
	return predicate.Organization(sql.FieldEQ(FieldName, v))
}

// Code applies equality check predicate on the "code" field. It's identical to CodeEQ.
func Code(v string) predicate.Organization {
	return predicate.Organization(sql.FieldEQ(FieldCode, v))
}

// LogoURL applies equality check predicate on the "logo_url" field. It's identical to LogoURLEQ.
func LogoURL(v string) predicate.Organization {
	return predicate.Organization(sql.FieldEQ(FieldLogoURL, v))
}

// Address applies equality check predicate on the "address" field. It's identical to AddressEQ.
func Address(v string) predicate.Organization {
	return predicate.Organization(sql.FieldEQ(FieldAddress, v))
}

// Phone applies equality check predicate on the "phone" field. It's identical to PhoneEQ.
func Phone(v string) predicate.Organization {
	return predicate.Organization(sql.FieldEQ(FieldPhone, v))
}

// Email applies equality check predicate on the "email" field. It's identical to EmailEQ.
func Email(v string) predicate.Organization {
	return predicate.Organization(sql.FieldEQ(FieldEmail, v))
}

// Website applies equality check predicate on the "website" field. It's identical to WebsiteEQ.
func Website(v string) predicate.Organization {
	return predicate.Organization(sql.FieldEQ(FieldWebsite, v))
}

// CreatedAt applies equality check predicate on the "created_at" field. It's identical to CreatedAtEQ.
func CreatedAt(v time.Time) predicate.Organization {
	return predicate.Organization(sql.FieldEQ(FieldCreatedAt, v))
}

// UpdatedAt applies equality check predicate on the "updated_at" field. It's identical to UpdatedAtEQ.
func UpdatedAt(v time.Time) predicate.Organization {
	return predicate.Organization(sql.FieldEQ(FieldUpdatedAt, v))
}

// ParentID applies equality check predicate on the "parent_id" field. It's identical to ParentIDEQ.
func ParentID(v int) predicate.Organization {
	return predicate.Organization(sql.FieldEQ(FieldParentID, v))
}

// NameEQ applies the EQ predicate on the "name" field.
func NameEQ(v string) predicate.Organization {
	return predicate.Organization(sql.FieldEQ(FieldName, v))
}

// NameNEQ applies the NEQ predicate on the "name" field.
func NameNEQ(v string) predicate.Organization {
	return predicate.Organization(sql.FieldNEQ(FieldName, v))
}

// NameIn applies the In predicate on the "name" field.
func NameIn(vs ...string) predicate.Organization {
	return predicate.Organization(sql.FieldIn(FieldName, vs...))
}

// NameNotIn applies the NotIn predicate on the "name" field.
func NameNotIn(vs ...string) predicate.Organization {
	return predicate.Organization(sql.FieldNotIn(FieldName, vs...))
}

// NameGT applies the GT predicate on the "name" field.
func NameGT(v string) predicate.Organization {
	return predicate.Organization(sql.FieldGT(FieldName, v))
}

// NameGTE applies the GTE predicate on the "name" field.
func NameGTE(v string) predicate.Organization {
	return predicate.Organization(sql.FieldGTE(FieldName, v))
}

// NameLT applies the LT predicate on the "name" field.
func NameLT(v string) predicate.Organization {
	return predicate.Organization(sql.FieldLT(FieldName, v))
}

// NameLTE applies the LTE predicate on the "name" field.
func NameLTE(v string) predicate.Organization {
	return predicate.Organization(sql.FieldLTE(FieldName, v))
}

// NameContains applies the Contains predicate on the "name" field.
func NameContains(v string) predicate.Organization {
	return predicate.Organization(sql.FieldContains(FieldName, v))
}

// NameHasPrefix applies the HasPrefix predicate on the "name" field.
func NameHasPrefix(v string) predicate.Organization {
	return predicate.Organization(sql.FieldHasPrefix(FieldName, v))
}

// NameHasSuffix applies the HasSuffix predicate on the "name" field.
func NameHasSuffix(v string) predicate.Organization {
	return predicate.Organization(sql.FieldHasSuffix(FieldName, v))
}

// NameEqualFold applies the EqualFold predicate on the "name" field.
func NameEqualFold(v string) predicate.Organization {
	return predicate.Organization(sql.FieldEqualFold(FieldName, v))
}

// NameContainsFold applies the ContainsFold predicate on the "name" field.
func NameContainsFold(v string) predicate.Organization {
	return predicate.Organization(sql.FieldContainsFold(FieldName, v))
}

// CodeEQ applies the EQ predicate on the "code" field.
func CodeEQ(v string) predicate.Organization {
	return predicate.Organization(sql.FieldEQ(FieldCode, v))
}

// CodeNEQ applies the NEQ predicate on the "code" field.
func CodeNEQ(v string) predicate.Organization {
	return predicate.Organization(sql.FieldNEQ(FieldCode, v))
}

// CodeIn applies the In predicate on the "code" field.
func CodeIn(vs ...string) predicate.Organization {
	return predicate.Organization(sql.FieldIn(FieldCode, vs...))
}

// CodeNotIn applies the NotIn predicate on the "code" field.
func CodeNotIn(vs ...string) predicate.Organization {
	return predicate.Organization(sql.FieldNotIn(FieldCode, vs...))
}

// CodeGT applies the GT predicate on the "code" field.
func CodeGT(v string) predicate.Organization {
	return predicate.Organization(sql.FieldGT(FieldCode, v))
}

// CodeGTE applies the GTE predicate on the "code" field.
func CodeGTE(v string) predicate.Organization {
	return predicate.Organization(sql.FieldGTE(FieldCode, v))
}

// CodeLT applies the LT predicate on the "code" field.
func CodeLT(v string) predicate.Organization {
	return predicate.Organization(sql.FieldLT(FieldCode, v))
}

// CodeLTE applies the LTE predicate on the "code" field.
func CodeLTE(v string) predicate.Organization {
	return predicate.Organization(sql.FieldLTE(FieldCode, v))
}

// CodeContains applies the Contains predicate on the "code" field.
func CodeContains(v string) predicate.Organization {
	return predicate.Organization(sql.FieldContains(FieldCode, v))
}

// CodeHasPrefix applies the HasPrefix predicate on the "code" field.
func CodeHasPrefix(v string) predicate.Organization {
	return predicate.Organization(sql.FieldHasPrefix(FieldCode, v))
}

// CodeHasSuffix applies the HasSuffix predicate on the "code" field.
func CodeHasSuffix(v string) predicate.Organization {
	return predicate.Organization(sql.FieldHasSuffix(FieldCode, v))
}

// CodeEqualFold applies the EqualFold predicate on the "code" field.
func CodeEqualFold(v string) predicate.Organization {
	return predicate.Organization(sql.FieldEqualFold(FieldCode, v))
}

// CodeContainsFold applies the ContainsFold predicate on the "code" field.
func CodeContainsFold(v string) predicate.Organization {
	return predicate.Organization(sql.FieldContainsFold(FieldCode, v))
}

// LogoURLEQ applies the EQ predicate on the "logo_url" field.
func LogoURLEQ(v string) predicate.Organization {
	return predicate.Organization(sql.FieldEQ(FieldLogoURL, v))
}

// LogoURLNEQ applies the NEQ predicate on the "logo_url" field.
func LogoURLNEQ(v string) predicate.Organization {
	return predicate.Organization(sql.FieldNEQ(FieldLogoURL, v))
}

// LogoURLIn applies the In predicate on the "logo_url" field.
func LogoURLIn(vs ...string) predicate.Organization {
	return predicate.Organization(sql.FieldIn(FieldLogoURL, vs...))
}

// LogoURLNotIn applies the NotIn predicate on the "logo_url" field.
func LogoURLNotIn(vs ...string) predicate.Organization {
	return predicate.Organization(sql.FieldNotIn(FieldLogoURL, vs...))
}

// LogoURLGT applies the GT predicate on the "logo_url" field.
func LogoURLGT(v string) predicate.Organization {
	return predicate.Organization(sql.FieldGT(FieldLogoURL, v))
}

// LogoURLGTE applies the GTE predicate on the "logo_url" field.
func LogoURLGTE(v string) predicate.Organization {
	return predicate.Organization(sql.FieldGTE(FieldLogoURL, v))
}

// LogoURLLT applies the LT predicate on the "logo_url" field.
func LogoURLLT(v string) predicate.Organization {
	return predicate.Organization(sql.FieldLT(FieldLogoURL, v))
}

// LogoURLLTE applies the LTE predicate on the "logo_url" field.
func LogoURLLTE(v string) predicate.Organization {
	return predicate.Organization(sql.FieldLTE(FieldLogoURL, v))
}

// LogoURLContains applies the Contains predicate on the "logo_url" field.
func LogoURLContains(v string) predicate.Organization {
	return predicate.Organization(sql.FieldContains(FieldLogoURL, v))
}

// LogoURLHasPrefix applies the HasPrefix predicate on the "logo_url" field.
func LogoURLHasPrefix(v string) predicate.Organization {
	return predicate.Organization(sql.FieldHasPrefix(FieldLogoURL, v))
}

// LogoURLHasSuffix applies the HasSuffix predicate on the "logo_url" field.
func LogoURLHasSuffix(v string) predicate.Organization {
	return predicate.Organization(sql.FieldHasSuffix(FieldLogoURL, v))
}

// LogoURLIsNil applies the IsNil predicate on the "logo_url" field.
func LogoURLIsNil() predicate.Organization {
	return predicate.Organization(sql.FieldIsNull(FieldLogoURL))
}

// LogoURLNotNil applies the NotNil predicate on the "logo_url" field.
func LogoURLNotNil() predicate.Organization {
	return predicate.Organization(sql.FieldNotNull(FieldLogoURL))
}

// LogoURLEqualFold applies the EqualFold predicate on the "logo_url" field.
func LogoURLEqualFold(v string) predicate.Organization {
	return predicate.Organization(sql.FieldEqualFold(FieldLogoURL, v))
}

// LogoURLContainsFold applies the ContainsFold predicate on the "logo_url" field.
func LogoURLContainsFold(v string) predicate.Organization {
	return predicate.Organization(sql.FieldContainsFold(FieldLogoURL, v))
}

// AddressEQ applies the EQ predicate on the "address" field.
func AddressEQ(v string) predicate.Organization {
	return predicate.Organization(sql.FieldEQ(FieldAddress, v))
}

// AddressNEQ applies the NEQ predicate on the "address" field.
func AddressNEQ(v string) predicate.Organization {
	return predicate.Organization(sql.FieldNEQ(FieldAddress, v))
}

// AddressIn applies the In predicate on the "address" field.
func AddressIn(vs ...string) predicate.Organization {
	return predicate.Organization(sql.FieldIn(FieldAddress, vs...))
}

// AddressNotIn applies the NotIn predicate on the "address" field.
func AddressNotIn(vs ...string) predicate.Organization {
	return predicate.Organization(sql.FieldNotIn(FieldAddress, vs...))
}

// AddressGT applies the GT predicate on the "address" field.
func AddressGT(v string) predicate.Organization {
	return predicate.Organization(sql.FieldGT(FieldAddress, v))
}

// AddressGTE applies the GTE predicate on the "address" field.
func AddressGTE(v string) predicate.Organization {
	return predicate.Organization(sql.FieldGTE(FieldAddress, v))
}

// AddressLT applies the LT predicate on the "address" field.
func AddressLT(v string) predicate.Organization {
	return predicate.Organization(sql.FieldLT(FieldAddress, v))
}

// AddressLTE applies the LTE predicate on the "address" field.
func AddressLTE(v string) predicate.Organization {
	return predicate.Organization(sql.FieldLTE(FieldAddress, v))
}

// AddressContains applies the Contains predicate on the "address" field.
func AddressContains(v string) predicate.Organization {
	return predicate.Organization(sql.FieldContains(FieldAddress, v))
}

// AddressHasPrefix applies the HasPrefix predicate on the "address" field.
func AddressHasPrefix(v string) predicate.Organization {
	return predicate.Organization(sql.FieldHasPrefix(FieldAddress, v))
}

// AddressHasSuffix applies the HasSuffix predicate on the "address" field.
func AddressHasSuffix(v string) predicate.Organization {
	return predicate.Organization(sql.FieldHasSuffix(FieldAddress, v))
}

// AddressIsNil applies the IsNil predicate on the "address" field.
func AddressIsNil() predicate.Organization {
	return predicate.Organization(sql.FieldIsNull(FieldAddress))
}

// AddressNotNil applies the NotNil predicate on the "address" field.
func AddressNotNil() predicate.Organization {
	return predicate.Organization(sql.FieldNotNull(FieldAddress))
}

// AddressEqualFold applies the EqualFold predicate on the "address" field.
func AddressEqualFold(v string) predicate.Organization {
	return predicate.Organization(sql.FieldEqualFold(FieldAddress, v))
}

// AddressContainsFold applies the ContainsFold predicate on the "address" field.
func AddressContainsFold(v string) predicate.Organization {
	return predicate.Organization(sql.FieldContainsFold(FieldAddress, v))
}

// PhoneEQ applies the EQ predicate on the "phone" field.
func PhoneEQ(v string) predicate.Organization {
	return predicate.Organization(sql.FieldEQ(FieldPhone, v))
}

// PhoneNEQ applies the NEQ predicate on the "phone" field.
func PhoneNEQ(v string) predicate.Organization {
	return predicate.Organization(sql.FieldNEQ(FieldPhone, v))
}

// PhoneIn applies the In predicate on the "phone" field.
func PhoneIn(vs ...string) predicate.Organization {
	return predicate.Organization(sql.FieldIn(FieldPhone, vs...))
}

// PhoneNotIn applies the NotIn predicate on the "phone" field.
func PhoneNotIn(vs ...string) predicate.Organization {
	return predicate.Organization(sql.FieldNotIn(FieldPhone, vs...))
}

// PhoneGT applies the GT predicate on the "phone" field.
func PhoneGT(v string) predicate.Organization {
	return predicate.Organization(sql.FieldGT(FieldPhone, v))
}

// PhoneGTE applies the GTE predicate on the "phone" field.
func PhoneGTE(v string) predicate.Organization {
	return predicate.Organization(sql.FieldGTE(FieldPhone, v))
}

// PhoneLT applies the LT predicate on the "phone" field.
func PhoneLT(v string) predicate.Organization {
	return predicate.Organization(sql.FieldLT(FieldPhone, v))
}

// PhoneLTE applies the LTE predicate on the "phone" field.
func PhoneLTE(v string) predicate.Organization {
	return predicate.Organization(sql.FieldLTE(FieldPhone, v))
}

// PhoneContains applies the Contains predicate on the "phone" field.
func PhoneContains(v string) predicate.Organization {
	return predicate.Organization(sql.FieldContains(FieldPhone, v))
}

// PhoneHasPrefix applies the HasPrefix predicate on the "phone" field.
func PhoneHasPrefix(v string) predicate.Organization {
	return predicate.Organization(sql.FieldHasPrefix(FieldPhone, v))
}

// PhoneHasSuffix applies the HasSuffix predicate on the "phone" field.
func PhoneHasSuffix(v string) predicate.Organization {
	return predicate.Organization(sql.FieldHasSuffix(FieldPhone, v))
}

// PhoneIsNil applies the IsNil predicate on the "phone" field.
func PhoneIsNil() predicate.Organization {
	return predicate.Organization(sql.FieldIsNull(FieldPhone))
}

// PhoneNotNil applies the NotNil predicate on the "phone" field.
func PhoneNotNil() predicate.Organization {
	return predicate.Organization(sql.FieldNotNull(FieldPhone))
}

// PhoneEqualFold applies the EqualFold predicate on the "phone" field.
func PhoneEqualFold(v string) predicate.Organization {
	return predicate.Organization(sql.FieldEqualFold(FieldPhone, v))
}

// PhoneContainsFold applies the ContainsFold predicate on the "phone" field.
func PhoneContainsFold(v string) predicate.Organization {
	return predicate.Organization(sql.FieldContainsFold(FieldPhone, v))
}

// EmailEQ applies the EQ predicate on the "email" field.
func EmailEQ(v string) predicate.Organization {
	return predicate.Organization(sql.FieldEQ(FieldEmail, v))
}

// EmailNEQ applies the NEQ predicate on the "email" field.
func EmailNEQ(v string) predicate.Organization {
	return predicate.Organization(sql.FieldNEQ(FieldEmail, v))
}

// EmailIn applies the In predicate on the "email" field.
func EmailIn(vs ...string) predicate.Organization {
	return predicate.Organization(sql.FieldIn(FieldEmail, vs...))
}

// EmailNotIn applies the NotIn predicate on the "email" field.
func EmailNotIn(vs ...string) predicate.Organization {
	return predicate.Organization(sql.FieldNotIn(FieldEmail, vs...))
}

// EmailGT applies the GT predicate on the "email" field.
func EmailGT(v string) predicate.Organization {
	return predicate.Organization(sql.FieldGT(FieldEmail, v))
}

// EmailGTE applies the GTE predicate on the "email" field.
func EmailGTE(v string) predicate.Organization {
	return predicate.Organization(sql.FieldGTE(FieldEmail, v))
}

// EmailLT applies the LT predicate on the "email" field.
func EmailLT(v string) predicate.Organization {
	return predicate.Organization(sql.FieldLT(FieldEmail, v))
}

// EmailLTE applies the LTE predicate on the "email" field.
func EmailLTE(v string) predicate.Organization {
	return predicate.Organization(sql.FieldLTE(FieldEmail, v))
}

// EmailContains applies the Contains predicate on the "email" field.
func EmailContains(v string) predicate.Organization {
	return predicate.Organization(sql.FieldContains(FieldEmail, v))
}

// EmailHasPrefix applies the HasPrefix predicate on the "email" field.
func EmailHasPrefix(v string) predicate.Organization {
	return predicate.Organization(sql.FieldHasPrefix(FieldEmail, v))
}

// EmailHasSuffix applies the HasSuffix predicate on the "email" field.
func EmailHasSuffix(v string) predicate.Organization {
	return predicate.Organization(sql.FieldHasSuffix(FieldEmail, v))
}

// EmailIsNil applies the IsNil predicate on the "email" field.
func EmailIsNil() predicate.Organization {
	return predicate.Organization(sql.FieldIsNull(FieldEmail))
}

// EmailNotNil applies the NotNil predicate on the "email" field.
func EmailNotNil() predicate.Organization {
	return predicate.Organization(sql.FieldNotNull(FieldEmail))
}

// EmailEqualFold applies the EqualFold predicate on the "email" field.
func EmailEqualFold(v string) predicate.Organization {
	return predicate.Organization(sql.FieldEqualFold(FieldEmail, v))
}

// EmailContainsFold applies the ContainsFold predicate on the "email" field.
func EmailContainsFold(v string) predicate.Organization {
	return predicate.Organization(sql.FieldContainsFold(FieldEmail, v))
}

// WebsiteEQ applies the EQ predicate on the "website" field.
func WebsiteEQ(v string) predicate.Organization {
	return predicate.Organization(sql.FieldEQ(FieldWebsite, v))
}

// WebsiteNEQ applies the NEQ predicate on the "website" field.
func WebsiteNEQ(v string) predicate.Organization {
	return predicate.Organization(sql.FieldNEQ(FieldWebsite, v))
}

// WebsiteIn applies the In predicate on the "website" field.
func WebsiteIn(vs ...string) predicate.Organization {
	return predicate.Organization(sql.FieldIn(FieldWebsite, vs...))
}

// WebsiteNotIn applies the NotIn predicate on the "website" field.
func WebsiteNotIn(vs ...string) predicate.Organization {
	return predicate.Organization(sql.FieldNotIn(FieldWebsite, vs...))
}

// WebsiteGT applies the GT predicate on the "website" field.
func WebsiteGT(v string) predicate.Organization {
	return predicate.Organization(sql.FieldGT(FieldWebsite, v))
}

// WebsiteGTE applies the GTE predicate on the "website" field.
func WebsiteGTE(v string) predicate.Organization {
	return predicate.Organization(sql.FieldGTE(FieldWebsite, v))
}

// WebsiteLT applies the LT predicate on the "website" field.
func WebsiteLT(v string) predicate.Organization {
	return predicate.Organization(sql.FieldLT(FieldWebsite, v))
}

// WebsiteLTE applies the LTE predicate on the "website" field.
func WebsiteLTE(v string) predicate.Organization {
	return predicate.Organization(sql.FieldLTE(FieldWebsite, v))
}

// WebsiteContains applies the Contains predicate on the "website" field.
func WebsiteContains(v string) predicate.Organization {
	return predicate.Organization(sql.FieldContains(FieldWebsite, v))
}

// WebsiteHasPrefix applies the HasPrefix predicate on the "website" field.
func WebsiteHasPrefix(v string) predicate.Organization {
	return predicate.Organization(sql.FieldHasPrefix(FieldWebsite, v))
}

// WebsiteHasSuffix applies the HasSuffix predicate on the "website" field.
func WebsiteHasSuffix(v string) predicate.Organization {
	return predicate.Organization(sql.FieldHasSuffix(FieldWebsite, v))
}

// WebsiteIsNil applies the IsNil predicate on the "website" field.
func WebsiteIsNil() predicate.Organization {
	return predicate.Organization(sql.FieldIsNull(FieldWebsite))
}

// WebsiteNotNil applies the NotNil predicate on the "website" field.
func WebsiteNotNil() predicate.Organization {
	return predicate.Organization(sql.FieldNotNull(FieldWebsite))
}

// WebsiteEqualFold applies the EqualFold predicate on the "website" field.
func WebsiteEqualFold(v string) predicate.Organization {
	return predicate.Organization(sql.FieldEqualFold(FieldWebsite, v))
}

// WebsiteContainsFold applies the ContainsFold predicate on the "website" field.
func WebsiteContainsFold(v string) predicate.Organization {
	return predicate.Organization(sql.FieldContainsFold(FieldWebsite, v))
}

// CreatedAtEQ applies the EQ predicate on the "created_at" field.
func CreatedAtEQ(v time.Time) predicate.Organization {
	return predicate.Organization(sql.FieldEQ(FieldCreatedAt, v))
}

// CreatedAtNEQ applies the NEQ predicate on the "created_at" field.
func CreatedAtNEQ(v time.Time) predicate.Organization {
	return predicate.Organization(sql.FieldNEQ(FieldCreatedAt, v))
}

// CreatedAtIn applies the In predicate on the "created_at" field.
func CreatedAtIn(vs ...time.Time) predicate.Organization {
	return predicate.Organization(sql.FieldIn(FieldCreatedAt, vs...))
}

// CreatedAtNotIn applies the NotIn predicate on the "created_at" field.
func CreatedAtNotIn(vs ...time.Time) predicate.Organization {
	return predicate.Organization(sql.FieldNotIn(FieldCreatedAt, vs...))
}

// CreatedAtGT applies the GT predicate on the "created_at" field.
func CreatedAtGT(v time.Time) predicate.Organization {
	return predicate.Organization(sql.FieldGT(FieldCreatedAt, v))
}

// CreatedAtGTE applies the GTE predicate on the "created_at" field.
func CreatedAtGTE(v time.Time) predicate.Organization {
	return predicate.Organization(sql.FieldGTE(FieldCreatedAt, v))
}

// CreatedAtLT applies the LT predicate on the "created_at" field.
func CreatedAtLT(v time.Time) predicate.Organization {
	return predicate.Organization(sql.FieldLT(FieldCreatedAt, v))
}

// CreatedAtLTE applies the LTE predicate on the "created_at" field.
func CreatedAtLTE(v time.Time) predicate.Organization {
	return predicate.Organization(sql.FieldLTE(FieldCreatedAt, v))
}

// UpdatedAtEQ applies the EQ predicate on the "updated_at" field.
func UpdatedAtEQ(v time.Time) predicate.Organization {
	return predicate.Organization(sql.FieldEQ(FieldUpdatedAt, v))
}

// UpdatedAtNEQ applies the NEQ predicate on the "updated_at" field.
func UpdatedAtNEQ(v time.Time) predicate.Organization {
	return predicate.Organization(sql.FieldNEQ(FieldUpdatedAt, v))
}

// UpdatedAtIn applies the In predicate on the "updated_at" field.
func UpdatedAtIn(vs ...time.Time) predicate.Organization {
	return predicate.Organization(sql.FieldIn(FieldUpdatedAt, vs...))
}

// UpdatedAtNotIn applies the NotIn predicate on the "updated_at" field.
func UpdatedAtNotIn(vs ...time.Time) predicate.Organization {
	return predicate.Organization(sql.FieldNotIn(FieldUpdatedAt, vs...))
}

// UpdatedAtGT applies the GT predicate on the "updated_at" field.
func UpdatedAtGT(v time.Time) predicate.Organization {
	return predicate.Organization(sql.FieldGT(FieldUpdatedAt, v))
}

// UpdatedAtGTE applies the GTE predicate on the "updated_at" field.
func UpdatedAtGTE(v time.Time) predicate.Organization {
	return predicate.Organization(sql.FieldGTE(FieldUpdatedAt, v))
}

// UpdatedAtLT applies the LT predicate on the "updated_at" field.
func UpdatedAtLT(v time.Time) predicate.Organization {
	return predicate.Organization(sql.FieldLT(FieldUpdatedAt, v))
}

// UpdatedAtLTE applies the LTE predicate on the "updated_at" field.
func UpdatedAtLTE(v time.Time) predicate.Organization {
	return predicate.Organization(sql.FieldLTE(FieldUpdatedAt, v))
}

// ParentIDEQ applies the EQ predicate on the "parent_id" field.
func ParentIDEQ(v int) predicate.Organization {
	return predicate.Organization(sql.FieldEQ(FieldParentID, v))
}

// ParentIDNEQ applies the NEQ predicate on the "parent_id" field.
func ParentIDNEQ(v int) predicate.Organization {
	return predicate.Organization(sql.FieldNEQ(FieldParentID, v))
}

// ParentIDIn applies the In predicate on the "parent_id" field.
func ParentIDIn(vs ...int) predicate.Organization {
	return predicate.Organization(sql.FieldIn(FieldParentID, vs...))
}

// ParentIDNotIn applies the NotIn predicate on the "parent_id" field.
func ParentIDNotIn(vs ...int) predicate.Organization {
	return predicate.Organization(sql.FieldNotIn(FieldParentID, vs...))
}

// ParentIDIsNil applies the IsNil predicate on the "parent_id" field.
func ParentIDIsNil() predicate.Organization {
	return predicate.Organization(sql.FieldIsNull(FieldParentID))
}

// ParentIDNotNil applies the NotNil predicate on the "parent_id" field.
func ParentIDNotNil() predicate.Organization {
	return predicate.Organization(sql.FieldNotNull(FieldParentID))
}

// HasParent applies the HasEdge predicate on the "parent" edge.
func HasParent() predicate.Organization {
	return predicate.Organization(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, ParentTable, ParentColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasParentWith applies the HasEdge predicate on the "parent" edge with a given conditions (other predicates).
func HasParentWith(preds ...predicate.Organization) predicate.Organization {
	return predicate.Organization(func(s *sql.Selector) {
		step := newParentStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasChildren applies the HasEdge predicate on the "children" edge.
func HasChildren() predicate.Organization {
	return predicate.Organization(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, ChildrenTable, ChildrenColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasChildrenWith applies the HasEdge predicate on the "children" edge with a given conditions (other predicates).
func HasChildrenWith(preds ...predicate.Organization) predicate.Organization {
	return predicate.Organization(func(s *sql.Selector) {
		step := newChildrenStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasDepartments applies the HasEdge predicate on the "departments" edge.
func HasDepartments() predicate.Organization {
	return predicate.Organization(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, DepartmentsTable, DepartmentsColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasDepartmentsWith applies the HasEdge predicate on the "departments" edge with a given conditions (other predicates).
func HasDepartmentsWith(preds ...predicate.Department) predicate.Organization {
	return predicate.Organization(func(s *sql.Selector) {
		step := newDepartmentsStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasProjects applies the HasEdge predicate on the "projects" edge.
func HasProjects() predicate.Organization {
	return predicate.Organization(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, ProjectsTable, ProjectsColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasProjectsWith applies the HasEdge predicate on the "projects" edge with a given conditions (other predicates).
func HasProjectsWith(preds ...predicate.Project) predicate.Organization {
	return predicate.Organization(func(s *sql.Selector) {
		step := newProjectsStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasLabels applies the HasEdge predicate on the "labels" edge.
func HasLabels() predicate.Organization {
	return predicate.Organization(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, LabelsTable, LabelsColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasLabelsWith applies the HasEdge predicate on the "labels" edge with a given conditions (other predicates).
func HasLabelsWith(preds ...predicate.Label) predicate.Organization {
	return predicate.Organization(func(s *sql.Selector) {
		step := newLabelsStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasLeaveRequests applies the HasEdge predicate on the "leave_requests" edge.
func HasLeaveRequests() predicate.Organization {
	return predicate.Organization(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, LeaveRequestsTable, LeaveRequestsColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasLeaveRequestsWith applies the HasEdge predicate on the "leave_requests" edge with a given conditions (other predicates).
func HasLeaveRequestsWith(preds ...predicate.LeaveRequest) predicate.Organization {
	return predicate.Organization(func(s *sql.Selector) {
		step := newLeaveRequestsStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.Organization) predicate.Organization {
	return predicate.Organization(sql.AndPredicates(predicates...))
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.Organization) predicate.Organization {
	return predicate.Organization(sql.OrPredicates(predicates...))
}

// Not applies the not operator on the given predicate.
func Not(p predicate.Organization) predicate.Organization {
	return predicate.Organization(sql.NotPredicates(p))
}
