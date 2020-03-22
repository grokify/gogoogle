package slidesutil

const (
	UnitUnspecified                  = "UNIT_UNSPECIFIED"              // The units are unknown.
	UnitEMU                          = "EMU"                           // An English Metric Unit (EMU) is defined as 1/360,000 of a centimeter
	UnitPT                           = "PT"                            // A point, 1/72 of an inch.
	LayoutUnspecified                = "PREDEFINED_LAYOUT_UNSPECIFIED" // Unspecified layout.
	LayoutBlank                      = "BLANK"                         // Blank layout, with no placeholders.
	LayoutCaptionOnly                = "CAPTION_ONLY"                  // Layout with a caption at the bottom.
	LayoutTitle                      = "TITLE"                         // Layout with a title and a subtitle.
	LayoutTitleAndBody               = "TITLE_AND_BODY"                // Layout with a title and body.
	LayoutTitleAndTwoColumns         = "TITLE_AND_TWO_COLUMNS"         // Layout with a title and two columns.
	LayoutTitleOnly                  = "TITLE_ONLY"                    // Layout with only a title.
	LayoutSectionHeader              = "SECTION_HEADER"                // Layout with a section title.
	LayoutSectionTitleAndDescription = "SECTION_TITLE_AND_DESCRIPTION" // Layout with a title and subtitle
	LayoutOneColumnText              = "ONE_COLUMN_TEXT"               // Layout with one title and one body, arranged in a single column.
	LayoutMainPoint                  = "MAIN_POINT"                    // Layout with a main point.
	LayoutBigNumber                  = "BIG_NUMBER"                    // Layout with a big number heading.
	RangeTypeUnspecified             = "RANGE_TYPE_UNSPECIFIED"
	RangeTypeFixedRange              = "FIXED_RANGE"
	RangeTypeStartIndex              = "FROM_START_INDEX"
	RangeTypeAll                     = "ALL"
)
