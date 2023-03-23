import React, { useEffect } from 'react'

import { TelemetryProps } from '@sourcegraph/shared/src/telemetry/telemetryService'

import { PageTitle } from '../../../../../components/PageTitle'

import { CodeInsightsExamples } from './components/code-insights-examples/CodeInsightsExamples'
import { CodeInsightsTemplates } from './components/code-insights-templates/CodeInsightsTemplates'
import { DynamicCodeInsightExample } from './components/dynamic-code-insight-example/DynamicCodeInsightExample'

import styles from './CodeInsightsGettingStartedPage.module.scss'

interface CodeInsightsGettingStartedPageProps extends TelemetryProps {
    isSourcegraphApp?: boolean
}

export const CodeInsightsGettingStartedPage: React.FunctionComponent<
    React.PropsWithChildren<CodeInsightsGettingStartedPageProps>
> = props => {
    const { telemetryService, isSourcegraphApp } = props

    useEffect(() => {
        telemetryService.logViewEvent('InsightsGetStartedPage')
    }, [telemetryService])

    return (
        <main className="pb-5">
            <PageTitle title="Code Insights" />
            <DynamicCodeInsightExample isSourcegraphApp={isSourcegraphApp} telemetryService={telemetryService} />
            <CodeInsightsExamples telemetryService={telemetryService} className={styles.section} />
            <CodeInsightsTemplates telemetryService={telemetryService} className={styles.section} />
        </main>
    )
}
