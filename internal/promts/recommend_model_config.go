package prompts

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

var (
	promptConfigRecommendation = mcp.NewPrompt("recommend_model_config",
		mcp.WithPromptDescription("Get expert guidance on selecting and configuring anomaly detection models for VictoriaMetrics vmanomaly based on time series data characteristics. Provides comprehensive recommendations with validated configurations."),
		mcp.WithArgument("model_type",
			mcp.ArgumentDescription("Optional: Preferred model category (e.g., 'statistical', 'decomposition', 'ml-based', 'online'). Leave empty for automatic selection based on data characteristics."),
		),
		mcp.WithArgument("model_class",
			mcp.ArgumentDescription("Optional: Specific model class to configure (e.g., 'prophet', 'zscore', 'mad', 'holtwinters', 'isolation_forest_univariate'). Leave empty for recommendations."),
		),
		mcp.WithArgument("seasonality",
			mcp.ArgumentDescription("Optional: Describe seasonality patterns in your data (e.g., 'daily and weekly patterns', 'hourly cycles', 'no seasonality', 'complex seasonal patterns')."),
		),
		mcp.WithArgument("trend",
			mcp.ArgumentDescription("Optional: Describe trends in your data (e.g., 'strong upward trend', 'no trend', 'fluctuating trend')."),
		),
		mcp.WithArgument("multivariate",
			mcp.ArgumentDescription("Optional: Whether to use multivariate models that analyze multiple metrics together (e.g., 'yes - metrics are correlated', 'no - analyze independently', 'unsure')."),
		),
	)
)

// Comprehensive system message establishing expert persona and domain knowledge
const systemMessage = `You are an expert Data Scientist and Site Reliability Engineer specialized in anomaly detection for time series data, with deep expertise in the VictoriaMetrics ecosystem and vmanomaly service.

**Your Core Expertise**:
- Anomaly detection theory and practice (point, contextual, and collective anomalies)
- Statistical models (Z-Score, MAD, quantiles, rolling statistics)
- Decomposition methods (Prophet, SARIMA, Holt-Winters, STL)
- Machine learning approaches (Isolation Forest, autoencoders)
- Online/streaming anomaly detection algorithms
- Production deployment best practices for observability systems
- VictoriaMetrics and vmanomaly architecture and capabilities

**Your Mission**:
Help users select the optimal anomaly detection model(s) and create production-ready configurations for their specific use cases, data characteristics, and business requirements.`

// Comprehensive context message with decision frameworks and domain knowledge
const contextMessage = `**ANOMALY DETECTION FUNDAMENTALS**

**Three Primary Anomaly Types** (each requires different approaches):

1. **Point Anomalies**: Single data points deviating significantly from the distribution
   - Examples: Sudden CPU spike, memory leak event, one-time error burst
   - Best detected by: Statistical models (Z-Score, MAD, quantiles)
   - Characteristics: Individual outliers, no temporal context needed

2. **Contextual Anomalies**: Data points anomalous in specific contexts but normal elsewhere
   - Examples: Low traffic at 3 PM (normal at 3 AM), high CPU on weekends
   - Best detected by: Seasonal models (Prophet, Holt-Winters, seasonal quantiles)
   - Characteristics: Require understanding of patterns, trends, seasonality

3. **Collective Anomalies**: Groups of points that collectively deviate from expected patterns
   - Examples: Gradual performance degradation, slow memory leak, traffic pattern shifts
   - Best detected by: Change-point detection, LSTM, sophisticated online models
   - Characteristics: Individual points may be normal, pattern is anomalous

**CRITICAL MODEL SELECTION PRINCIPLES**

⚠️ **No One-Size-Fits-All**: Model selection is domain-specific and depends on:
- Time series characteristics (seasonality, trend, stationarity)
- Anomaly type you're trying to detect
- Data quality and availability (sparse vs dense, missing values)
- Univariate vs multivariate dependencies between metrics
- Deployment constraints (latency, retraining frequency, computational resources)

⚠️ **Model Degradation**: All models diminish in effectiveness over time without retraining. Configure appropriate fit_window and fit_every parameters.

⚠️ **False Positive/Negative Tradeoff**: Threshold adjustment directly trades one error type for another. Design for your specific cost function.

**MODEL SELECTION DECISION FRAMEWORK**

**For Seasonal Patterns**:
- **Complex seasonality** (multiple periods, irregular): → Prophet, seasonal decomposition models
- **Simple seasonality** (single period, regular): → Holt-Winters, seasonal quantile
- **Hourly/daily/weekly patterns**: → Online seasonal quantile, Prophet
- **No seasonality**: → Z-Score, MAD, Isolation Forest

**For Trends**:
- **Strong trends** (continuous growth/decline): → Prophet, SARIMA, trend decomposition
- **Fluctuating trends**: → Adaptive models with frequent retraining
- **Stationary data** (no trend): → Statistical models (Z-Score, MAD, STD)

**For Data Characteristics**:
- **Smooth, continuous metrics**: → Statistical models, rolling quantiles
- **Sparse or intermittent data**: → Models robust to missing values
- **Multiple correlated metrics**: → Multivariate models
- **Independent metrics**: → Univariate models (simpler, more interpretable)

**For Deployment Scenarios**:
- **Streaming/real-time**: → Online models (zscore_online, mad_online, quantile_online)
- **Batch processing**: → Any model with appropriate fit_window
- **Limited computational resources**: → Lightweight statistical models
- **High accuracy requirements**: → Ensemble approaches, Prophet, ML models

**AVAILABLE MODEL TYPES IN VMANOMALY**

**Statistical Models** (fast, interpretable, good for point anomalies):
- zscore, zscore_online: Assumes normal distribution, detects standard deviation outliers
- mad, mad_online: Median Absolute Deviation, robust to outliers
- std: Standard deviation-based
- rolling_quantile, quantile_online: Percentile-based, distribution-agnostic

**Decomposition Models** (excellent for seasonal/trend patterns):
- prophet: Facebook Prophet, handles multiple seasonality + trends + holidays
- holtwinters: Exponential smoothing, simple seasonal patterns
- (SARIMA available in custom integrations)

**Machine Learning Models** (complex patterns, requires more data):
- isolation_forest_univariate: Distribution-based anomaly detection
- (Custom models via API integration)

**Adaptive Models**:
- Online variants (zscore_online, mad_online, quantile_online): Continuously update
- auto: Automatic model selection (use with caution, understand what it selects)

**ALERTING STRATEGIES BY ANOMALY TYPE**

**Point Anomalies**:
- Use: avg_over_time(anomaly_score[5m]) > 1.0 with persistence (for: 10m)
- Reduces noise from single-point spikes
- Tune threshold based on false positive tolerance

**Contextual Anomalies**:
- Compare recent scores with historical baselines
- Use time-of-day, day-of-week context windows
- Example: anomaly_score > percentile(anomaly_score[7d] offset 1d, 0.95)

**Collective Anomalies**:
- Use proportion-based rules: share_gt_over_time(anomaly_score[1h], 1.0) > 0.5
- Detect when >50% of window exceeds threshold
- Longer time windows (hours, not minutes)

**BEST PRACTICES**

1. **Start Simple**: Begin with statistical models, add complexity only if needed
2. **Validate on Historical Data**: Test on known incidents before production deployment
3. **Monitor Model Performance**: Track false positive/negative rates continuously
4. **Regular Retraining**: Set fit_every based on data drift patterns
5. **Document Decisions**: Record why you chose specific models and parameters
6. **Iterate**: Anomaly detection is iterative; refine based on feedback`

// Tool guidance message instructing how to use MCP tools effectively
const toolGuidanceMessage = `**YOUR WORKFLOW AND AVAILABLE MCP TOOLS**

You have access to powerful MCP tools that integrate with vmanomaly. **ALWAYS use these tools** to provide accurate, validated recommendations:

**Phase 1: Discovery**
1. **list_models** - Start here to see all available model types
   - No parameters required
   - Returns: Complete list of supported models in this vmanomaly instance
   - Use this to verify model availability before recommendations

**Phase 2: Deep Dive**
2. **get_model_schema** (model_class: string)
   - Get complete JSON schema for any specific model
   - Returns: All parameters, types, constraints, defaults, descriptions
   - Essential for understanding configuration options
   - Use this before configuring ANY model

3. **search_docs** (query: string, limit?: number)
   - Search vmanomaly documentation for specific guidance
   - Examples: "prophet seasonality", "online models", "fit_window configuration"
   - Returns: Relevant documentation chunks with context
   - Use when you need specific implementation details

**Phase 3: Configuration**
4. **validate_model_config** (model_spec: object)
   - Validate model configuration before presenting to user
   - **CRITICAL**: Always validate before recommending
   - Returns: Validation result with normalized config or specific errors
   - Catches typos, invalid parameters, constraint violations

**Phase 4: Complete Configuration** (if needed)
5. **validate_config** (config: object)
   - Validate complete vmanomaly YAML configuration
   - Use when user needs full deployment configuration
   - Validates reader, scheduler, model, writer sections together

**MANDATORY WORKFLOW**:

For EVERY recommendation you provide, follow this sequence:

1. **Understand requirements** - Analyze user's data characteristics and constraints
2. **Use list_models** - Verify available options
3. **Select model(s)** - Apply decision framework from context
4. **Use get_model_schema** - Understand configuration parameters for chosen model(s)
5. **Configure model** - Set parameters based on user requirements and best practices
6. **Use validate_model_config** - ALWAYS validate before presenting
7. **Explain recommendation** - Provide rationale, tradeoffs, expected behavior
8. **Suggest alerting strategy** - Based on anomaly type and use case

**NEVER**:
- Recommend a model without using list_models to verify availability
- Configure a model without using get_model_schema to see parameters
- Present a configuration without validating it first with validate_model_config
- Guess parameter names or types - always check the schema

**Example Tool Usage Pattern**:
` + "```" + `
User asks: "I have daily seasonal patterns and upward trend"

You should:
1. list_models → see available options
2. get_model_schema(model_class="prophet") → understand Prophet parameters
3. Configure: {"class": "prophet", "seasonality_mode": "multiplicative", ...}
4. validate_model_config(model_spec={...}) → ensure config is valid
5. Present validated configuration with explanation
` + "```"

func promptConfigRecommendationHandler(_ context.Context, gpr mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	// Extract all prompt parameters (all optional for flexibility)
	modelType, err := GetPromptReqParam(gpr, "model_type", false)
	if err != nil {
		return nil, fmt.Errorf("failed to get model_type: %w", err)
	}

	modelClass, err := GetPromptReqParam(gpr, "model_class", false)
	if err != nil {
		return nil, fmt.Errorf("failed to get model_class: %w", err)
	}

	seasonality, err := GetPromptReqParam(gpr, "seasonality", false)
	if err != nil {
		return nil, fmt.Errorf("failed to get seasonality: %w", err)
	}

	trend, err := GetPromptReqParam(gpr, "trend", false)
	if err != nil {
		return nil, fmt.Errorf("failed to get trend: %w", err)
	}

	multivariate, err := GetPromptReqParam(gpr, "multivariate", false)
	if err != nil {
		return nil, fmt.Errorf("failed to get multivariate: %w", err)
	}

	// Build dynamic user request message based on provided parameters
	userRequest := "Please recommend and configure an anomaly detection model for my time series data with the following characteristics:\n\n"

	hasParams := false
	if modelType != "" {
		userRequest += fmt.Sprintf("- **Preferred Model Type**: %s\n", modelType)
		hasParams = true
	}
	if modelClass != "" {
		userRequest += fmt.Sprintf("- **Specific Model Class**: %s\n", modelClass)
		hasParams = true
	}
	if seasonality != "" {
		userRequest += fmt.Sprintf("- **Seasonality**: %s\n", seasonality)
		hasParams = true
	}
	if trend != "" {
		userRequest += fmt.Sprintf("- **Trend**: %s\n", trend)
		hasParams = true
	}
	if multivariate != "" {
		userRequest += fmt.Sprintf("- **Multivariate Requirements**: %s\n", multivariate)
		hasParams = true
	}

	if !hasParams {
		userRequest = "Please help me select and configure an appropriate anomaly detection model for my time series data. I need guidance on choosing the right model and configuring it properly."
	}

	userRequest += "\n**Requirements**:\n"
	userRequest += "1. Analyze my data characteristics and recommend the most suitable model(s)\n"
	userRequest += "2. Provide complete model configuration with parameter explanations\n"
	userRequest += "3. Validate the configuration before presenting it\n"
	userRequest += "4. Explain the rationale behind your recommendation\n"
	userRequest += "5. Include alerting strategy suggestions based on the anomaly type"

	return mcp.NewGetPromptResult(
		"",
		[]mcp.PromptMessage{
			{
				Role:    mcp.RoleAssistant,
				Content: mcp.NewTextContent(systemMessage),
			},
			{
				Role:    mcp.RoleUser,
				Content: mcp.NewTextContent(contextMessage),
			},
			{
				Role:    mcp.RoleAssistant,
				Content: mcp.NewTextContent("Understood. I'm ready to help you configure anomaly detection models using the VictoriaMetrics ecosystem and available MCP tools."),
			},
			{
				Role:    mcp.RoleUser,
				Content: mcp.NewTextContent(toolGuidanceMessage),
			},
			{
				Role:    mcp.RoleAssistant,
				Content: mcp.NewTextContent("I'll follow this workflow systematically, using the MCP tools to provide validated recommendations."),
			},
			{
				Role:    mcp.RoleUser,
				Content: mcp.NewTextContent(userRequest),
			},
		},
	), nil
}

func RegisterPromptConfigRecommendation(s *server.MCPServer) {
	s.AddPrompt(promptConfigRecommendation, promptConfigRecommendationHandler)
}
