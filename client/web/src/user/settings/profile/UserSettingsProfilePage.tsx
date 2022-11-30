import React, { useEffect } from 'react'

import { gql } from '@sourcegraph/http-client'
import { PageHeader, Link, Text } from '@sourcegraph/wildcard'

import { PageTitle } from '../../../components/PageTitle'
import { Timestamp } from '../../../components/time/Timestamp'
import { EditUserProfilePage as EditUserProfilePageFragment } from '../../../graphql-operations'
import { eventLogger } from '../../../tracking/eventLogger'

import { EditUserProfileForm } from './EditUserProfileForm'

import styles from './UserSettingsProfilePage.module.scss'

export const EditUserProfilePageGQLFragment = gql`
    fragment EditUserProfilePage on User {
        id
        username
        displayName
        avatarURL
        viewerCanChangeUsername
        createdAt
    }
`

interface Props {
    user: EditUserProfilePageFragment
}

export const UserSettingsProfilePage: React.FunctionComponent<React.PropsWithChildren<Props>> = ({
    user,
    ...props
}) => {
    useEffect(() => eventLogger.logViewEvent('UserProfile'), [])

    return (
        <div>
            <PageTitle title="Profile" />
            <PageHeader
                path={[{ text: 'Profile' }]}
                headingElement="h2"
                description={
                    <>
                        {user.displayName ? (
                            <>
                                {user.displayName} ({user.username})
                            </>
                        ) : (
                            user.username
                        )}{' '}
                        started using Sourcegraph <Timestamp date={user.createdAt} />.
                    </>
                }
                className={styles.heading}
            />
            {user && (
                <EditUserProfileForm
                    user={user}
                    initialValue={user}
                    after={
                        window.context.sourcegraphDotComMode && (
                            <Text className="mt-4">
                                <Link to="https://about.sourcegraph.com/contact">Contact support</Link> to delete your
                                account.
                            </Text>
                        )
                    }
                />
            )}
        </div>
    )
}
